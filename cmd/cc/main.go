package main

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Сообщения для таймеров
type (
	tick1Msg struct{}
	tick2Msg struct{}
)

type model struct {
	ctx       context.Context    // Контекст для отмены
	cancel    context.CancelFunc // Функция отмены
	timer1Cmd tea.Cmd            // Команда первого таймера
	timer2Cmd tea.Cmd            // Команда второго таймера
	count1    int                // Счётчик первого таймера
	count2    int                // Счётчик второго таймера
}

// Создаёт новую модель с чистым контекстом
func newModel() model {
	m := model{
		count1: 0,
		count2: 0,
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())
	return m
}

func (m model) Init() tea.Cmd {
	// Запускаем оба таймера при инициализации
	return tea.Batch(m.startTimer1(), m.startTimer2())
}

// Таймер на 1 секунду
func (m model) startTimer1() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		select {
		case <-m.ctx.Done(): // Проверяем отмену контекста
			return nil
		default:
			return tick1Msg{}
		}
	})
}

// Таймер на 2 секунды
func (m model) startTimer2() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		select {
		case <-m.ctx.Done(): // Проверяем отмену контекста
			return nil
		default:
			return tick2Msg{}
		}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "r" {
			// 1. Отменяем текущие таймеры
			m.cancel()
			// 2. Создаём совершенно новую модель
			newModel := newModel()
			// 3. Запускаем её таймеры
			return newModel, newModel.Init()
		}
		if msg.Type == tea.KeyCtrlC {
			m.cancel()
			return m, tea.Quit
		}

	case tick1Msg:
		m.count1++
		return m, m.startTimer1()

	case tick2Msg:
		m.count2++
		return m, m.startTimer2()
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf(
		"Таймер 1 (1 сек): %d\n"+
			"Таймер 2 (2 сек): %d\n\n"+
			"Нажмите 'r' для сброса\n"+
			"Ctrl+C для выхода",
		m.count1, m.count2,
	)
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ошибка: %v", err)
	}
}
