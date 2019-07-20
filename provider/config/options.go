package config

type Options struct {
	Sources []string
}

type Option func(*Options)
