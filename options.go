package pathmatch

type Option func(p *Path) error

func SetSeperator(sep string) Option {
	return func(p *Path) error {
		p.Seperator = sep
		return nil
	}
}

func SetWildcard(wildcard string) Option {
	return func(p *Path) error {
		p.Wildcard = wildcard
		return nil
	}
}

func SetPrefix(prefix string) Option {
	return func(p *Path) error {
		p.Prefix = prefix
		return nil
	}
}

func SetSuffix(suffix string) Option {
	return func(p *Path) error {
		p.Suffix = suffix
		return nil
	}
}
