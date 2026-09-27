package gui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	omarchytheme "diary/internal/adapter/theme"
	"diary/internal/domain"
	"diary/internal/usecase"
)

// DesktopApp represents the native desktop GUI application.
type DesktopApp struct {
	fyneApp      fyne.App
	window       fyne.Window
	useCase      usecase.DailyGoalUseCase
	themeService omarchytheme.Service

	currentDate domain.Date
	currentGoal *domain.DailyGoal

	// UI Widgets
	dateLabel     *widget.Label
	progressLabel *widget.Label
	progressBar   *widget.ProgressBar
	newGoalEntry  *widget.Entry
	tasksBox      *fyne.Container
	statusLabel   *widget.Label
}

// NewDesktopApp initializes the native desktop application.
func NewDesktopApp(useCase usecase.DailyGoalUseCase, themeService omarchytheme.Service) *DesktopApp {
	if themeService == nil {
		themeService = omarchytheme.NewOmarchyService("")
	}

	a := app.NewWithID("com.omarchy.diary")
	currentTheme := themeService.GetCurrentTheme()
	a.Settings().SetTheme(NewOmarchyTheme(currentTheme))

	w := a.NewWindow("Diary - Daily Goals")
	w.Resize(fyne.NewSize(480, 620))

	if len(AppIconBytes) > 0 {
		w.SetIcon(fyne.NewStaticResource("icon.png", AppIconBytes))
	}

	desktop := &DesktopApp{
		fyneApp:      a,
		window:       w,
		useCase:      useCase,
		themeService: themeService,
		currentDate:  domain.Today(),
	}

	desktop.setupUI()
	desktop.startThemeWatcher()

	return desktop
}

// Run starts the native desktop application event loop.
func (d *DesktopApp) Run() {
	d.loadDate(d.currentDate)
	d.window.ShowAndRun()
}

func (d *DesktopApp) setupUI() {
	// 1. Branding Header
	title := widget.NewLabelWithStyle("Diary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("Daily Goals")

	var brandIcon fyne.CanvasObject
	if len(AppIconBytes) > 0 {
		img := canvas.NewImageFromResource(fyne.NewStaticResource("icon.png", AppIconBytes))
		img.SetMinSize(fyne.NewSize(26, 26))
		brandIcon = img
	} else {
		brandIcon = widget.NewIcon(theme.DocumentIcon())
	}

	brandRow := container.NewHBox(brandIcon, title, subtitle)

	// 2. Date Navigation
	prevBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		d.loadDate(d.currentDate.AddDays(-1))
	})
	nextBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		d.loadDate(d.currentDate.AddDays(1))
	})
	todayBtn := widget.NewButton("Today", func() {
		d.loadDate(domain.Today())
	})

	d.dateLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	dateNav := container.NewBorder(nil, nil, prevBtn, container.NewHBox(nextBtn, todayBtn), d.dateLabel)

	// 3. Progress Module
	d.progressLabel = widget.NewLabel("0 of 0 completed (0%)")
	d.progressBar = widget.NewProgressBar()
	d.progressBar.SetValue(0)
	progressBox := container.NewVBox(d.progressLabel, d.progressBar)

	// 4. Input Row
	d.newGoalEntry = widget.NewEntry()
	d.newGoalEntry.SetPlaceHolder("What do you want to accomplish today? (Press Enter)")
	d.newGoalEntry.OnSubmitted = func(text string) {
		d.addGoal(text)
	}

	addBtn := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		d.addGoal(d.newGoalEntry.Text)
	})

	inputRow := container.NewBorder(nil, nil, nil, addBtn, d.newGoalEntry)

	// 5. Tasks Stack Container
	d.tasksBox = container.NewVBox()
	tasksScroll := container.NewVScroll(d.tasksBox)

	// 6. Footer / Status
	d.statusLabel = widget.NewLabelWithStyle("Saved locally in Markdown", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	// Layout Composition
	topBox := container.NewVBox(
		brandRow,
		widget.NewSeparator(),
		dateNav,
		progressBox,
		inputRow,
		widget.NewSeparator(),
	)

	content := container.NewBorder(
		topBox,
		container.NewVBox(widget.NewSeparator(), d.statusLabel),
		nil,
		nil,
		tasksScroll,
	)

	// Shortcuts
	d.window.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyLeft:
			d.loadDate(d.currentDate.AddDays(-1))
		case fyne.KeyRight:
			d.loadDate(d.currentDate.AddDays(1))
		}
	})

	d.window.SetContent(container.NewPadded(content))
}

func (d *DesktopApp) loadDate(target domain.Date) {
	d.currentDate = target
	d.dateLabel.SetText(d.formatDateHeading(target))

	goal, err := d.useCase.GetDailyGoals(context.Background(), target.String())
	if err != nil {
		d.statusLabel.SetText(fmt.Sprintf("Error loading goals: %v", err))
		return
	}

	d.currentGoal = goal
	d.renderGoals()
}

func (d *DesktopApp) renderGoals() {
	d.tasksBox.Objects = nil

	items := d.currentGoal.Items()
	total := len(items)
	completed := d.currentGoal.DoneCount()

	// Update Progress
	var fraction float64
	if total > 0 {
		fraction = float64(completed) / float64(total)
	}
	d.progressBar.SetValue(fraction)
	d.progressLabel.SetText(fmt.Sprintf("%d of %d completed (%d%%)", completed, total, d.currentGoal.ProgressPercentage()))

	if total == 0 {
		emptyMsg := widget.NewLabelWithStyle("No goals for this day yet.\nAdd a new goal above to get started!", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
		d.tasksBox.Add(container.NewCenter(emptyMsg))
		d.tasksBox.Refresh()
		d.statusLabel.SetText(fmt.Sprintf("File: ./metas/%s.md", d.currentDate.String()))
		return
	}

	for i, item := range items {
		idx := i
		it := item

		// Checkbox
		chk := widget.NewCheck("", func(checked bool) {
			d.toggleGoal(idx)
		})
		chk.Checked = it.Done()

		// Text label with visual style
		var label *widget.Label
		if it.Done() {
			label = widget.NewLabelWithStyle(it.Text(), fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
		} else {
			label = widget.NewLabelWithStyle(it.Text(), fyne.TextAlignLeading, fyne.TextStyle{Bold: false})
		}

		// Delete button
		delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			d.deleteGoal(idx)
		})

		row := container.NewBorder(nil, nil, chk, delBtn, label)
		d.tasksBox.Add(row)
	}

	d.tasksBox.Refresh()
	d.statusLabel.SetText(fmt.Sprintf("File: ./metas/%s.md (saved)", d.currentDate.String()))
}

func (d *DesktopApp) addGoal(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	newItem, err := domain.NewItem(text, false)
	if err != nil {
		d.statusLabel.SetText(fmt.Sprintf("Invalid text: %v", err))
		return
	}
	newItems := append(d.currentGoal.Items(), newItem)

	saved, err := d.useCase.SaveDailyGoals(context.Background(), d.currentDate.String(), newItems)
	if err != nil {
		d.statusLabel.SetText(fmt.Sprintf("Error saving: %v", err))
		return
	}

	d.currentGoal = saved
	d.newGoalEntry.SetText("")
	d.renderGoals()
}

func (d *DesktopApp) toggleGoal(idx int) {
	items := d.currentGoal.Items()
	if idx < 0 || idx >= len(items) {
		return
	}

	newItems := make([]domain.Item, len(items))
	copy(newItems, items)
	newItems[idx] = newItems[idx].Toggle()

	saved, err := d.useCase.SaveDailyGoals(context.Background(), d.currentDate.String(), newItems)
	if err != nil {
		d.statusLabel.SetText(fmt.Sprintf("Error toggling goal: %v", err))
		return
	}

	d.currentGoal = saved
	d.renderGoals()
}

func (d *DesktopApp) deleteGoal(idx int) {
	items := d.currentGoal.Items()
	if idx < 0 || idx >= len(items) {
		return
	}

	newItems := append([]domain.Item{}, items[:idx]...)
	newItems = append(newItems, items[idx+1:]...)

	saved, err := d.useCase.SaveDailyGoals(context.Background(), d.currentDate.String(), newItems)
	if err != nil {
		d.statusLabel.SetText(fmt.Sprintf("Error deleting: %v", err))
		return
	}

	d.currentGoal = saved
	d.renderGoals()
}

func (d *DesktopApp) formatDateHeading(dt domain.Date) string {
	t := dt.Time()
	return t.Format("Monday, 02 January 2006")
}

func (d *DesktopApp) startThemeWatcher() {
	go func() {
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()

		var lastThemeName string
		for range ticker.C {
			info := d.themeService.GetCurrentTheme()
			if info.ThemeName != lastThemeName {
				lastThemeName = info.ThemeName
				d.fyneApp.Settings().SetTheme(NewOmarchyTheme(info))
			}
		}
	}()
}
