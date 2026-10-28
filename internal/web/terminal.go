package web

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tclasen/factoryctl/internal/interaction"
	"github.com/tclasen/factoryctl/internal/project"
)

// ChatTerminal is a human answer surface, never the MCP protocol's stdin.
func ChatTerminal(projects []project.Project, cfg RuntimeConfig, prompt string, in io.Reader, out io.Writer) error {
	file, ok := in.(*os.File)
	if !ok {
		return fmt.Errorf("interaction_unavailable: chat requires a terminal")
	}
	if !isTerminal(file) {
		return fmt.Errorf("interaction_unavailable: chat requires a terminal")
	}
	if len(projects) != 1 {
		return fmt.Errorf("select one project for chat")
	}
	h, e := New(projects)
	if e != nil {
		return e
	}
	h.EnableConversations(cfg)
	defer h.Close()
	listener, e := Listen("127.0.0.1:0")
	if e != nil {
		return e
	}
	defer listener.Close()
	h.chat.base = "http://" + listener.Addr().String()
	server := &http.Server{Handler: h, ReadHeaderTimeout: 5 * time.Second, WriteTimeout: 45 * time.Second}
	defer server.Close()
	go server.Serve(listener)
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 4096), 16384)
	if strings.TrimSpace(prompt) == "" {
		fmt.Fprint(out, "Your message: ")
		if !scanner.Scan() {
			return fmt.Errorf("interaction_unavailable: no initial input")
		}
		prompt = scanner.Text()
	}
	scope, token, e := h.chat.broker.NewScope(projects[0].ID)
	if e != nil {
		return e
	}
	chat := &conversation{ID: randomToken(), Project: projects[0].ID, Title: prompt, Status: "starting", Messages: []chatMessage{{Role: "user", Text: prompt}}, Native: []nativeRequest{}, scope: scope, token: token, root: projects[0].Root}
	h.chat.mu.Lock()
	h.chat.conversations[chat.ID] = chat
	h.chat.mu.Unlock()
	go h.chat.start(chat, projects[0], prompt)
	fmt.Fprintln(out, "FactoryCTL conversation. Type /quit to close, /cancel to interrupt. Generated views use plain text here.")
	lines := make(chan string)
	ended := make(chan struct{})
	go func() {
		defer close(ended)
		for scanner.Scan() {
			select {
			case lines <- scanner.Text():
			case <-h.chat.ctx.Done():
				return
			}
		}
	}()
	ticker := time.NewTicker(150 * time.Millisecond)
	defer ticker.Stop()
	shown := map[string]bool{}
	printed := map[string]string{}
	lastStatus := ""
	for {
		select {
		case <-ended:
			return fmt.Errorf("interaction_unavailable: terminal input closed")
		case <-ticker.C:
			h.chat.mu.Lock()
			messages := append([]chatMessage(nil), chat.Messages...)
			status, errText := chat.Status, chat.Error
			native := append([]nativeRequest(nil), chat.Native...)
			h.chat.mu.Unlock()
			for _, m := range messages {
				if m.Role == "assistant" && m.ID != "" && printed[m.ID] != m.Text {
					if strings.HasPrefix(m.Text, printed[m.ID]) {
						fmt.Fprint(out, strings.TrimPrefix(m.Text, printed[m.ID]))
					} else {
						fmt.Fprintln(out, m.Text)
					}
					printed[m.ID] = m.Text
				}
			}
			snap, _ := h.chat.broker.Snapshot(scope)
			for _, q := range snap.Questions {
				if q.Status == "pending" && !shown[q.ID] {
					fmt.Fprintf(out, "\nQuestion [%s]: %s\n", q.ID, q.Prompt)
					for _, choice := range q.Choices {
						fmt.Fprintf(out, "  %s: %s\n", choice.ID, choice.Label)
					}
					fmt.Fprintln(out, "Enter your answer (choice IDs separated by commas, or custom text); /decline refuses.")
					shown[q.ID] = true
				}
			}
			for _, v := range snap.Views {
				if !shown[v.ID] {
					fmt.Fprintf(out, "\nView: %s\n%s\n", v.Title, v.Fallback)
					shown[v.ID] = true
				}
			}
			for _, q := range native {
				if !shown[q.ID] {
					fmt.Fprintf(out, "\nNative runtime request %s: %s\n%s\nRespond with JSON matching this request (no automatic approval).\n", q.ID, q.Method, q.Params)
					shown[q.ID] = true
				}
			}
			if status != lastStatus {
				if status == "idle" {
					fmt.Fprint(out, "\nYour follow-up: ")
				}
				if status == "failed" || status == "disconnected" {
					return fmt.Errorf("%s: %s", status, errText)
				}
				lastStatus = status
			}
		case line := <-lines:
			if line == "/quit" {
				return nil
			}
			h.chat.mu.Lock()
			client, thread, turn, status := chat.client, chat.ThreadID, chat.TurnID, chat.Status
			native := append([]nativeRequest(nil), chat.Native...)
			h.chat.mu.Unlock()
			if line == "/cancel" {
				h.chat.mu.Lock()
				client, thread, turn = chat.client, chat.ThreadID, chat.TurnID
				active := client != nil && turn != ""
				if active {
					chat.Status = "interrupting"
					h.chat.broker.CancelPending(scope, "cancelled")
				}
				h.chat.mu.Unlock()
				if active {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					_, e = client.Request(ctx, "turn/interrupt", map[string]string{"threadId": thread, "turnId": turn})
					cancel()
					if e != nil {
						fmt.Fprintln(out, e)
					}
				}
				continue
			}
			if len(native) > 0 {
				q := native[0]
				raw := json.RawMessage(line)
				if e := validateNativeAnswer(q, raw); e != nil {
					fmt.Fprintln(out, e)
					continue
				}
				if e := client.Reply(json.RawMessage(q.ID), raw); e != nil {
					fmt.Fprintln(out, e)
					continue
				}
				h.chat.mu.Lock()
				if len(chat.Native) > 0 && chat.Native[0].ID == q.ID {
					chat.Native = chat.Native[1:]
				}
				h.chat.mu.Unlock()
				continue
			}
			snap, _ := h.chat.broker.Snapshot(scope)
			var pending *interaction.Question
			for i := range snap.Questions {
				if snap.Questions[i].Status == "pending" {
					pending = &snap.Questions[i]
					break
				}
			}
			if pending != nil {
				answer := interaction.AnswerInput{Status: "answered", Text: line}
				if line == "/decline" {
					answer.Status = "declined"
					answer.Text = ""
				} else if pending.Kind != "text" {
					ids := strings.Split(line, ",")
					all := true
					for i, id := range ids {
						ids[i] = strings.TrimSpace(id)
						found := false
						for _, choice := range pending.Choices {
							if choice.ID == ids[i] {
								found = true
							}
						}
						all = all && found
					}
					if all {
						answer.OptionIDs = ids
						answer.Text = ""
					}
				}
				if _, e := h.chat.broker.Answer(scope, pending.ID, answer); e != nil {
					fmt.Fprintln(out, e)
				}
				continue
			}
			if status != "idle" && status != "interrupted" {
				fmt.Fprintln(out, "Wait for the current turn or use /cancel.")
				continue
			}
			if strings.TrimSpace(line) == "" {
				continue
			}
			h.chat.mu.Lock()
			chat.Status = "running"
			chat.Messages = append(chat.Messages, chatMessage{Role: "user", Text: line})
			h.chat.mu.Unlock()
			go h.chat.turn(chat, line)
		}
	}
}
