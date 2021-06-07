package main

import (
	lg "github.com/charmbracelet/lipgloss"
)

var (
	statusBarLogoStyle       = lg.NewStyle().Bold(true).Background(lg.Color("208")).Foreground(lg.Color("255"))
	statusBarNameStyle       = lg.NewStyle().Bold(true).Background(lg.Color("208")).Foreground(lg.Color("16"))
	tabStyle                 = lg.NewStyle().Background(lg.Color("8")).Foreground(lg.Color("15"))
	selectorTitleStyle       = lg.NewStyle().Foreground(lg.Color("15"))
	selectorTitleActiveStyle = lg.NewStyle().Foreground(lg.Color("3"))
	selectorUrlStyle         = lg.NewStyle().Foreground(lg.Color("8"))
	selectorNumStyle         = lg.NewStyle().Foreground(lg.Color("7"))
	storyCommentUserStyle    = lg.NewStyle().Foreground(lg.Color("15"))
	storyCommentTimeStyle    = lg.NewStyle().Underline(true)
	storyCommentKidsStyle    = lg.NewStyle().BorderStyle(lg.RoundedBorder()).BorderLeft(true).PaddingLeft(1)
	storyTitleStyle          = lg.NewStyle().Foreground(lg.Color("15")).Bold(true)
)
