package jobs

import "context"

func (s *Service) HandleDeleteOldCommands(ctx context.Context) error {
	return s.commandService.DeleteOldCommands(ctx)
}
