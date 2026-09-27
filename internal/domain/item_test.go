package domain_test

import (
	"testing"

	"diary/internal/domain"
)

func TestNewItem(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		done      bool
		wantText  string
		wantDone  bool
		wantError error
	}{
		{
			name:      "valid item",
			text:      "Estudar Clean Architecture",
			done:      false,
			wantText:  "Estudar Clean Architecture",
			wantDone:  false,
			wantError: nil,
		},
		{
			name:      "item with newlines and extra spaces",
			text:      "  Primeira linha\nSegunda linha\r  ",
			done:      true,
			wantText:  "Primeira linha Segunda linha",
			wantDone:  true,
			wantError: nil,
		},
		{
			name:      "empty item text",
			text:      "   \n\r  ",
			done:      false,
			wantText:  "",
			wantDone:  false,
			wantError: domain.ErrEmptyItemText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := domain.NewItem(tt.text, tt.done)
			if tt.wantError != nil {
				if err != tt.wantError {
					t.Fatalf("expected error %v, got %v", tt.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if item.Text() != tt.wantText {
				t.Errorf("got text %q, want %q", item.Text(), tt.wantText)
			}
			if item.Done() != tt.wantDone {
				t.Errorf("got done %v, want %v", item.Done(), tt.wantDone)
			}
		})
	}
}

func TestItem_Toggle(t *testing.T) {
	item, err := domain.NewItem("Tarefa", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	toggled := item.Toggle()
	if !toggled.Done() {
		t.Errorf("expected item to be done")
	}

	untoggled := toggled.Toggle()
	if untoggled.Done() {
		t.Errorf("expected item to be undone")
	}
}
