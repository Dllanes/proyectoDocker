package mongo

import (
	"github.com/guadalupej/proyecto/pkg/models"
	"github.com/guadalupej/proyecto/pkg/newerrors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (r *Repository) InsertEpisode(episode models.Episode) error {
	session := r.Session.Copy()
	defer session.Close()
	com := session.DB(r.DatabaseName).C("episodes")
	err := com.Insert(&episode)
	if err != nil {
		if mgo.IsDup(err) {
			return newerrors.NewErrBadRequest("id already exists")
		}
		return err
	}
	return nil
}

func (r *Repository) GetEpisodes(filters models.EpisodesFilters) ([]models.Episode, error) {
	session := r.Session.Copy()
	defer session.Close()
	com := session.DB(r.DatabaseName).C("episodes")
	searchQuery := bson.M{}
	if filters.Name != "" {
		searchQuery["name"] = filters.Name
	}
	if filters.Episode != "" {
		searchQuery["episode"] = filters.Episode
	}
	episodes := []models.Episode{}
	err := com.Find(searchQuery).Limit(filters.Limit).Skip(filters.Offset).All(&episodes)
	if err != nil {
		return nil, err
	}
	return episodes, nil
}

func (r *Repository) GetEpisodeByID(id string) (*models.Episode, error) {
	session := r.Session.Copy()
	defer session.Close()
	com := session.DB(r.DatabaseName).C("episodes")
	episode := models.Episode{}
	err := com.FindId(id).One(&episode)
	if err != nil {
		if err.Error() == "not found" {
			return nil, newerrors.NewErrNotFound("episode not found")
		}
		return nil, err
	}

	return &episode, nil
}
