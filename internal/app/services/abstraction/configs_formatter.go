package abstraction

type ConfigsFormatter interface {
	Format(rawConfig string) (string, error)
}
