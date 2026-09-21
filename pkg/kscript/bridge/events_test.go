package bridge

import (
	"context"
	"runtime"
	"strings"
	"sync"
	"testing"

	taskrun "github.com/gookit/taskrun"

	"github.com/inhere/kite-go/pkg/kscript"
)

// eventRecorder collects library events delivered through the bridge.
type eventRecorder struct {
	mu     sync.Mutex
	events []taskrun.Event
}

func (r *eventRecorder) observe(event taskrun.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
}

func (r *eventRecorder) all() []taskrun.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]taskrun.Event(nil), r.events...)
}

// TestBridgeDeliversEvents checks the observer seam Kite binds for its verbose
// output: the host renders progress, the library only reports data.
func TestBridgeDeliversEvents(t *testing.T) {
	rec := &eventRecorder{}
	legacy := newLegacy(t, map[string]any{
		"tpl": map[string]any{
			"vars": map[string]any{"who": "events"},
			"run":  "echo ${who}",
		},
	})
	var out strings.Builder
	b := New(legacy,
		WithBaseDir(t.TempDir()),
		WithEventObserver(rec.observe),
		WithIO(taskrun.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}),
	)
	found, err := b.TryRun(context.Background(), "tpl", nil, &kscript.RunCtx{Silent: true})
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	events := rec.all()
	if len(events) == 0 {
		t.Fatal("no events were delivered")
	}
	kinds := map[taskrun.EventKind]int{}
	for _, event := range events {
		kinds[event.Kind]++
		if event.Task != "tpl" || event.CallID == "" {
			t.Fatalf("event lacks coordinates: %+v", event)
		}
	}
	if kinds[taskrun.EventRunStarted] != 1 || kinds[taskrun.EventRunFinished] != 1 {
		t.Fatalf("run events=%v", kinds)
	}
	if kinds[taskrun.EventStepStarted] != 1 || kinds[taskrun.EventStepFinished] != 1 {
		t.Fatalf("step events=%v", kinds)
	}
	for _, event := range events {
		if event.Kind == taskrun.EventRunFinished && event.Status != taskrun.StatusSucceeded {
			t.Fatalf("run finished event=%+v", event)
		}
	}
}

// TestBridgeReportsSkippedTaskEvent covers a converted task that is skipped by
// its condition: Kite sees a skip reason instead of a started event.
func TestBridgeReportsSkippedTaskEvent(t *testing.T) {
	rec := &eventRecorder{}
	legacy := newLegacy(t, map[string]any{
		"caller": map[string]any{"run": "@task:off"},
		"off": map[string]any{
			"if":  "false",
			"run": "echo never",
		},
	})
	var out strings.Builder
	b := New(legacy,
		WithBaseDir(t.TempDir()),
		WithEventObserver(rec.observe),
		WithIO(taskrun.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}),
	)
	if _, err := b.TryRun(context.Background(), "caller", nil, &kscript.RunCtx{Silent: true}); err != nil {
		t.Fatalf("TryRun: %v", err)
	}
	var skipped *taskrun.Event
	var started bool
	for _, event := range rec.all() {
		switch event.Kind {
		case taskrun.EventTaskSkipped:
			if event.Task == "off" {
				copied := event
				skipped = &copied
			}
		case taskrun.EventTaskStarted:
			if event.Task == "off" {
				started = true
			}
		}
	}
	if skipped == nil || skipped.Reason == "" {
		t.Fatalf("no skip event with a reason: %+v", rec.all())
	}
	if started {
		t.Fatal("a skipped task must not report started")
	}
}

// TestBridgeDefaultObserverLogs checks the default observer forwards to the Kite
// logger without changing the run result.
func TestBridgeDefaultObserverLogs(t *testing.T) {
	shell := "sh"
	if runtime.GOOS == "windows" {
		shell = "cmd"
	}
	script := "printf logged"
	if runtime.GOOS == "windows" {
		script = "echo logged"
	}
	legacy := newLegacy(t, map[string]any{
		"log": map[string]any{"type": shell, "run": script},
	})
	var out strings.Builder
	b := New(legacy, WithBaseDir(t.TempDir()), WithIO(taskrun.IO{Stdout: &out, Stderr: &out, CaptureLimit: 1 << 16}))
	if found, err := b.TryRun(context.Background(), "log", nil, &kscript.RunCtx{Silent: true}); err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if !strings.Contains(out.String(), "logged") {
		t.Fatalf("output=%q", out.String())
	}
}
