package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/xfiendx4life/gb_back_2_hw/hw9/pkg/models"
)

type Local struct {
	Path string
}

// mode - "n" for new dir, empty for existing dir
func New(path string, mode string) (Storage, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("can't read filepath: %s", err)
	}
	_, err = os.ReadDir(absPath)
	if err != nil {
		if mode == "" {
			return nil, fmt.Errorf("can't open dir: %s", err)
		}
		err = os.Mkdir(absPath, 0777)
		if err != nil {
			return nil, fmt.Errorf("can't create dir: %s", err)
		}
		_, err = os.ReadDir(absPath)
		if err != nil {
			return nil, fmt.Errorf("can't read created dir: %s", err)
		}
	}
	return &Local{
		Path: absPath,
	}, nil
}

func (l *Local) Create(ctx context.Context, list models.List) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		path := fmt.Sprintf("%s/%s", l.Path, list.ID.String())
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("can't create file: %s", err)
		}
		defer file.Close()
		log.Printf("file %s created", path)
		data, err := json.Marshal(list)
		if err != nil {
			return fmt.Errorf("can't marshal list to json: %s", err)
		}
		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf("can't write json to file: %s", err)
		}
		log.Printf("data to file %s is written", path)
		return nil
	}
}

func (l *Local) Read(ctx context.Context, id uuid.UUID) (list *models.List, err error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("done with context")
	default:
		data, err := os.ReadFile(filepath.Join(l.Path, id.String()))
		if err != nil {
			return nil, fmt.Errorf("can't read data %s", err)
		}
		log.Printf("got data %s", string(data))
		err = json.Unmarshal(data, &list)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal data: %s", err)
		}
		return list, err
	}
}

func findItem(name string, data []*models.Item) (item *models.Item, ok bool) {
	for _, itm := range data {
		if itm.Name == name {
			return itm, true
		}
	}
	return nil, false
}

// Update doesn't delete any items from price, just updates prices
func (l *Local) Update(ctx context.Context, id uuid.UUID, newItems []*models.Item) (error) {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		lst, err := l.Read(ctx, id)
		if err != nil {
			return fmt.Errorf("can't read file")
		}
		log.Printf("Read data from file with id %s\n", id.String())
		for _, itm := range newItems {
			if found, ok := findItem(itm.Name, lst.Items); ok {
				found.Price = itm.Price
				log.Printf("change price of item with name %s to %d\n", found.Name, itm.Price)
			}
		}
		data, err := json.Marshal(lst)
		if err != nil {
			return fmt.Errorf("can't marshal json %s", err)
		}
		log.Printf("Marshal new data to json\n")
		p := filepath.Join(l.Path, id.String())
		err = os.WriteFile(p, data, 0777)
		if err != nil {
			return fmt.Errorf("can't write to file %s", err)
		}
		log.Printf("Wrote to file %s new data\n", p)
		return nil
	}
}
func (l *Local) Delete(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		err := os.Remove(filepath.Join(l.Path, id.String()))
		if err != nil {
			return fmt.Errorf("can't delete list: %s", err)
		}
		return nil
	}
}
