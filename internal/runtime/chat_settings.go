package runtime

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tclasen/factoryctl/internal/safeio"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

const maxChatSettings = 64 << 10

var ErrSettingsConflict = errors.New("chat settings changed; reload before saving")

type ChatSelection struct {
	Model  string `json:"model"`
	Effort string `json:"reasoning_effort"`
}
type ChatSettingsSnapshot struct {
	Global   *ChatSelection           `json:"global"`
	Projects map[string]ChatSelection `json:"projects"`
	Revision string                   `json:"revision"`
}
type chatSettingsDocument struct {
	Version  int                      `json:"version"`
	Global   *ChatSelection           `json:"global"`
	Projects map[string]ChatSelection `json:"projects"`
}
type ChatSettingsStore struct{ path string }

func (s *ChatSettingsStore) Path() string { return s.path }

func NewChatSettingsStore(path string) (*ChatSettingsStore, error) {
	if path == "" {
		u, err := user.Current()
		if err != nil {
			return nil, errors.New("cannot resolve current user's settings home")
		}
		path = filepath.Join(u.HomeDir, ".factoryctl", "chat-settings.json")
	}
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return nil, errors.New("chat settings path must be absolute and clean")
	}
	return &ChatSettingsStore{path: path}, nil
}
func CanonicalProjectRoot(path string) (string, error) {
	if !filepath.IsAbs(path) {
		return "", errors.New("project root must be absolute")
	}
	p, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	st, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	if !st.IsDir() {
		return "", errors.New("project root must be a directory")
	}
	return filepath.Clean(p), nil
}
func ResolveChatSelection(explicit *ChatSelection, invocation ChatSelection, saved ChatSettingsSnapshot, root string, legacy ChatSelection) ChatSelection {
	if explicit != nil {
		return *explicit
	}
	selected := legacy
	if saved.Global != nil {
		selected = *saved.Global
	}
	if project, ok := saved.Projects[root]; ok {
		selected = project
	}
	if invocation.Model != "" {
		return invocation
	}
	if invocation.Effort != "" {
		selected.Effort = invocation.Effort
	}
	return selected
}
func validChatSelection(s ChatSelection) error {
	for _, v := range []string{s.Model, s.Effort} {
		if !utf8.ValidString(v) || len(v) > 512 || strings.ContainsAny(v, "\x00\r\n") || strings.TrimSpace(v) != v {
			return errors.New("invalid chat model or reasoning value")
		}
	}
	return nil
}
func emptyChatSettings() ChatSettingsSnapshot {
	return ChatSettingsSnapshot{Projects: map[string]ChatSelection{}, Revision: "absent"}
}
func (s *ChatSettingsStore) Read() (ChatSettingsSnapshot, error) {
	root, err := s.directory(false)
	if errors.Is(err, os.ErrNotExist) {
		return emptyChatSettings(), nil
	}
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	defer root.Close()
	return s.readRoot(root)
}
func (s *ChatSettingsStore) readRoot(root *os.Root) (ChatSettingsSnapshot, error) {
	b, err := safeio.ReadRoot(root, filepath.Base(s.path), maxChatSettings)
	if errors.Is(err, os.ErrNotExist) {
		return emptyChatSettings(), nil
	}
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	if len(b) > maxChatSettings || !utf8.Valid(b) {
		return ChatSettingsSnapshot{}, errors.New("chat settings exceed bounds or contain invalid UTF-8")
	}
	if err = uniqueJSON(b); err != nil {
		return ChatSettingsSnapshot{}, err
	}
	if err = validateStoredPairs(b); err != nil {
		return ChatSettingsSnapshot{}, err
	}
	var doc chatSettingsDocument
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err = d.Decode(&doc); err != nil {
		return ChatSettingsSnapshot{}, fmt.Errorf("invalid chat settings: %w", err)
	}
	if doc.Version != 1 || len(doc.Projects) > 256 {
		return ChatSettingsSnapshot{}, errors.New("unsupported chat settings version or project count")
	}
	if doc.Global != nil {
		if err = validChatSelection(*doc.Global); err != nil {
			return ChatSettingsSnapshot{}, err
		}
	}
	if doc.Projects == nil {
		doc.Projects = map[string]ChatSelection{}
	}
	for path, selection := range doc.Projects {
		if len(path) > 4096 || !filepath.IsAbs(path) || filepath.Clean(path) != path || strings.ContainsRune(path, 0) {
			return ChatSettingsSnapshot{}, errors.New("invalid saved project root")
		}
		if err = validChatSelection(selection); err != nil {
			return ChatSettingsSnapshot{}, err
		}
	}
	sum := sha256.Sum256(b)
	return ChatSettingsSnapshot{Global: doc.Global, Projects: doc.Projects, Revision: hex.EncodeToString(sum[:])}, nil
}

// directory refuses a symlink at the settings directory itself. Ancestors are
// resolved by the OS (e.g. macOS /var); subsequent operations use a rooted handle.
func (s *ChatSettingsStore) directory(create bool) (*os.Root, error) {
	parent := filepath.Dir(s.path)
	st, err := os.Lstat(parent)
	if errors.Is(err, os.ErrNotExist) && create {
		if err = os.Mkdir(parent, 0700); err != nil && !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		st, err = os.Lstat(parent)
	}
	if err != nil {
		return nil, err
	}
	if !st.IsDir() || st.Mode()&os.ModeSymlink != 0 {
		return nil, errors.New("unsafe chat settings directory")
	}
	root, err := os.OpenRoot(parent)
	if err != nil {
		return nil, err
	}
	current, err := root.Stat(".")
	if err != nil || !os.SameFile(st, current) {
		root.Close()
		return nil, errors.New("chat settings directory changed")
	}
	return root, nil
}
func (s *ChatSettingsStore) Save(ctx context.Context, projectRoot string, selection *ChatSelection, expectedRevision string) (ChatSettingsSnapshot, error) {
	if expectedRevision == "" {
		return ChatSettingsSnapshot{}, ErrSettingsConflict
	}
	if projectRoot != "" {
		var err error
		projectRoot, err = CanonicalProjectRoot(projectRoot)
		if err != nil {
			return ChatSettingsSnapshot{}, err
		}
	}
	if selection != nil {
		if err := validChatSelection(*selection); err != nil {
			return ChatSettingsSnapshot{}, err
		}
	}
	root, err := s.directory(true)
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	defer root.Close()
	lockName := filepath.Base(s.path) + ".lock"
	lock, err := root.OpenFile(lockName, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if errors.Is(err, os.ErrExist) {
		lock, err = safeio.OpenRegularRoot(root, lockName, os.O_RDWR)
	}
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	defer lock.Close()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	for {
		if err = ctx.Err(); err != nil {
			return ChatSettingsSnapshot{}, err
		}
		err = lockChatSettings(lock)
		if err == nil {
			break
		}
		if !chatLockBusy(err) {
			return ChatSettingsSnapshot{}, err
		}
		select {
		case <-ctx.Done():
			return ChatSettingsSnapshot{}, ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
	defer unlockChatSettings(lock)
	held, err := lock.Stat()
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	present, err := root.Lstat(lockName)
	if err != nil || !os.SameFile(held, present) || !present.Mode().IsRegular() {
		return ChatSettingsSnapshot{}, ErrSettingsConflict
	}
	current, err := s.readRoot(root)
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	if current.Revision != expectedRevision {
		return ChatSettingsSnapshot{}, ErrSettingsConflict
	}
	if projectRoot == "" {
		current.Global = selection
	} else if selection == nil {
		delete(current.Projects, projectRoot)
	} else {
		current.Projects[projectRoot] = *selection
	}
	if len(current.Projects) > 256 {
		return ChatSettingsSnapshot{}, errors.New("too many saved projects")
	}
	b, err := json.MarshalIndent(chatSettingsDocument{Version: 1, Global: current.Global, Projects: current.Projects}, "", "  ")
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	b = append(b, '\n')
	if len(b) > maxChatSettings {
		return ChatSettingsSnapshot{}, errors.New("chat settings exceed size limit")
	}
	var nonce [16]byte
	if _, err = rand.Read(nonce[:]); err != nil {
		return ChatSettingsSnapshot{}, err
	}
	name := ".chat-settings-" + hex.EncodeToString(nonce[:]) + ".tmp"
	f, err := root.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	defer root.Remove(name)
	_, err = f.Write(b)
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	// Check again immediately before replacement; rooted rename replaces, never follows, a target.
	latest, err := s.readRoot(root)
	if err != nil {
		return ChatSettingsSnapshot{}, err
	}
	if latest.Revision != expectedRevision {
		return ChatSettingsSnapshot{}, ErrSettingsConflict
	}
	if err = root.Rename(name, filepath.Base(s.path)); err != nil {
		return ChatSettingsSnapshot{}, err
	}
	if dir, e := safeio.OpenDirectoryRoot(root, "."); e == nil {
		_ = dir.Sync()
		_ = dir.Close()
	}
	return s.readRoot(root)
}
func uniqueJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))
	var value func(int) error
	value = func(depth int) error {
		if depth > 8 {
			return errors.New("chat settings JSON nesting exceeds limit")
		}
		t, e := d.Token()
		if e != nil {
			return e
		}
		delim, ok := t.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				key, e := d.Token()
				if e != nil {
					return e
				}
				k, ok := key.(string)
				if !ok || seen[k] {
					return errors.New("duplicate chat settings field")
				}
				seen[k] = true
				if e = value(depth + 1); e != nil {
					return e
				}
			}
		case '[':
			for d.More() {
				if e = value(depth + 1); e != nil {
					return e
				}
			}
		default:
			return errors.New("invalid chat settings JSON")
		}
		_, e = d.Token()
		return e
	}
	if err := value(0); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return errors.New("trailing chat settings JSON")
	}
	return nil
}

func validateStoredPairs(b []byte) error {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(b, &top); err != nil || top == nil {
		return errors.New("chat settings must be an object")
	}
	pair := func(raw json.RawMessage) error {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(raw, &fields); err != nil || len(fields) != 2 {
			return errors.New("stored chat selection requires model and reasoning_effort strings")
		}
		for _, key := range []string{"model", "reasoning_effort"} {
			value, ok := fields[key]
			var text string
			if !ok || bytes.Equal(bytes.TrimSpace(value), []byte("null")) || json.Unmarshal(value, &text) != nil {
				return errors.New("stored chat selection requires model and reasoning_effort strings")
			}
		}
		return nil
	}
	if raw, ok := top["global"]; ok && !bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		if err := pair(raw); err != nil {
			return err
		}
	}
	if raw, ok := top["projects"]; ok {
		var projects map[string]json.RawMessage
		if err := json.Unmarshal(raw, &projects); err != nil || projects == nil {
			return errors.New("stored projects must be an object")
		}
		for _, value := range projects {
			if err := pair(value); err != nil {
				return err
			}
		}
	}
	return nil
}
