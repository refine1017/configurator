package models

import (
	"errors"
	"fmt"
	"sort"
)

type ConfigActivity struct {
	env       string
	Timestamp string `json:"timestamp"`
	Color     string `json:"color"`
	Content   string `json:"content"`
}

func MergeConfigInfo(env *Environment, collect string) ([]*ConfigActivity, string, error) {
	var activities ActivitySorter
	var mergeError = ""

	if env.Copy == "" {
		return activities, "Not a copy, cannot merge", nil
	}

	oldEnv, err := GetEnvironmentById(env.Copy)
	if err != nil {
		return activities, "", err
	}

	oldConfig := oldEnv.GetConfigByName(collect)
	if oldConfig == nil {
		return nil, "", errors.New("oldConfig not found")
	}

	newConfig := env.GetConfigByName(collect)
	if newConfig == nil {
		return nil, "", errors.New("newConfig not found")
	}

	activities = append(activities, &ConfigActivity{
		env:       oldEnv.Name,
		Timestamp: oldConfig.Updated,
		Content:   fmt.Sprintf("[%v] %v: %v", oldEnv.Name, oldConfig.Editor, oldConfig.Log),
		Color:     "#0bbd87",
	})

	activities = append(activities, &ConfigActivity{
		env:       env.Name,
		Timestamp: newConfig.Created,
		Content:   fmt.Sprintf("[%v] %v: copy from %v", env.Name, newConfig.Creator, oldEnv.Name),
		Color:     "#0bbd87",
	})

	if newConfig.Updated != newConfig.Created {
		activities = append(activities, &ConfigActivity{
			env:       env.Name,
			Timestamp: newConfig.Updated,
			Content:   fmt.Sprintf("[%v] %v: %v", env.Name, newConfig.Editor, newConfig.Log),
			Color:     "#0bbd87",
		})
	} else {
		mergeError = "No change, don't need merge"
	}

	sort.Sort(activities)

	for i, act := range activities {
		if i == 0 {
			continue
		}

		if act.env == oldEnv.Name {
			act.Color = "#ff0000"
			mergeError = "Find conflict, need to merge manually"
		}
	}

	return activities, mergeError, nil
}

func MergeConfig(env *Environment, collect string, admin string) error {
	_, mergeError, err := MergeConfigInfo(env, collect)
	if err != nil {
		return err
	}

	if mergeError != "" {
		return errors.New("cannot merge")
	}

	oldEnv, err := GetEnvironmentById(env.Copy)
	if err != nil {
		return err
	}

	oldConfig := oldEnv.GetConfigByName(collect)
	if oldConfig == nil {
		return errors.New("oldConfig not found")
	}

	newConfig := env.GetConfigByName(collect)
	if newConfig == nil {
		return errors.New("newConfig not found")
	}

	if err := configCol(oldEnv.Database, collect).DropCollection(); err != nil {
		return err
	}

	switch oldConfig.Format {
	case FormatTable:
		if err = CopyConfigTable(env, collect, oldEnv, collect); err != nil {
			return fmt.Errorf("CopyConfigTable:%v with err: %v", collect, err)
		}
	case FormatKV:
		if err = CopyConfigKV(env, collect, oldEnv, collect); err != nil {
			return fmt.Errorf("CopyConfigKV:%v with err: %v", collect, err)
		}
	case FormatJson:
		if err = CopyConfigJson(env, collect, oldEnv, collect); err != nil {
			return fmt.Errorf("CopyConfigJson:%v with err: %v", collect, err)
		}
	}

	oldConfig.SetEditor(admin, fmt.Sprintf("merge from %v", env.Name))
	if err := oldEnv.Save(admin); err != nil {
		return err
	}

	newConfig.SetCreator(admin, fmt.Sprintf("merge to %v", env.Name))
	if err := env.Save(admin); err != nil {
		return err
	}

	return nil
}

type ActivitySorter []*ConfigActivity

func (s ActivitySorter) Len() int      { return len(s) }
func (s ActivitySorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ActivitySorter) Less(i, j int) bool {
	return s[i].Timestamp < s[j].Timestamp
}
