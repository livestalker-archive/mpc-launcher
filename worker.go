package main

const (
	PlayCmd  = "887"
	PauseCmd = "888"
)

type Workers []Worker

type Worker struct {
	Args []string
	Port string
}
