package abstraction

type ConfigsOrganizer interface {
	Organize(groupId int, domainsCountLimit int) error
}
