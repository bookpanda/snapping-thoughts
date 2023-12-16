package seed

import "github.com/rs/zerolog/log"

func (s Seed) CourseSeed12180464() error {
	for _, i := range items {
		item, err := s.db.GetItemWithId(i.Id)
		if err != nil {
			log.Error().Str("seed", "Error with GetItemWithId: "+i.Id).Err(err)
			return err
		}

		if item.IsUsed != "" {
			i.IsUsed = item.IsUsed
		}
		err = s.db.CreateItem(i)

		if err != nil {
			return err
		}
	}
	return nil
}
