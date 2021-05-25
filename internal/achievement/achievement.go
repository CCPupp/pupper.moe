package achievement

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Achis struct {
	Achis []Achi `json:"Achis"`
}

type Achi struct {
	Id                 int    `json:"id"`
	Stage              int    `json:"stage"`
	AccuracyStage      int    `json:"accuracy_stage"`
	PrecisionStage     int    `json:"precision_stage"`
	ReadingStage       int    `json:"reading_stage"`
	SpeedStage         int    `json:"speed_stage"`
	StaminaStage       int    `json:"stamina_stage"`
	AccuracyStageNext  string `json:"accuracy_stage_next"`
	PrecisionStageNext string `json:"precision_stage_next"`
	ReadingStageNext   string `json:"reading_stage_next"`
	SpeedStageNext     string `json:"speed_stage_next"`
	StaminaStageNext   string `json:"stamina_stage_next"`
}

type Score struct {
	Id        int      `json:"id"`
	BestId    int      `json:"best_id"`
	UserId    int      `json:"user_id"`
	Accuracy  float64  `json:"accuracy"`
	Mods      []string `json:"mods"`
	Score     int      `json:"score"`
	MaxCombo  int      `json:"max_combo"`
	Perfect   bool     `json:"perfect"`
	Pp        float64  `json:"pp"`
	Rank      string   `json:"rank"`
	CreatedAt string   `json:"created_at"`
	Mode      string   `json:"mode"`
	ModeInt   int      `json:"int"`
	Replay    bool     `json:"replay"`
	Beatmap   Beatmap  `json:"beatmap"`
}

type Beatmap struct {
	URL string `json:"url"`
}

func CheckCompletion(recent []Score) {
	if len(recent) > 0 {
		achi := GetAchi(recent[0].UserId)

		for i := 0; i < len(recent); i++ {
			if achi.Stage == 0 {
				if recent[i].Beatmap.URL == "https://osu.ppy.sh/beatmaps/75" {
					//Stage 0 -> 1
					setTutorialDone(recent[i].UserId)
				}
			}
			if achi.Stage == 1 {
				if recent[i].Beatmap.URL != "" && recent[i].Accuracy > 0.99 {
					//Accuracy 1 -> 2
					setStage(achi, 2, "acc", "SS Any Map > 500 combo")
				}
				if recent[i].Beatmap.URL == "" {
					//Precision 1 -> 2
					setStage(achi, 2, "prec", "WIP")
				}
				if recent[i].Beatmap.URL == "" {
					//Reading 1 -> 2
					setStage(achi, 2, "read", "WIP")
				}
				if recent[i].Beatmap.URL == "" {
					//Speed 1 -> 2
					setStage(achi, 2, "speed", "WIP")
				}
				if recent[i].Beatmap.URL == "" {
					//Stamina 1 -> 2
					setStage(achi, 2, "stam", "WIP")
				}
				if recent[i].Beatmap.URL == "" {
					//Total 1 -> 2
					setStage(achi, 2, "total", "WIP")
				}
			}
		}
	}
}

func setStage(achi Achi, stage int, stageName, stageNext string) {

	currentList := GetAchiJSON()

	for i := 0; i < len(currentList.Achis); i++ {
		if currentList.Achis[i].Id == achi.Id {
			if stageName == "total" {
				currentList.Achis[i].Stage = stage
			} else if stageName == "acc" {
				currentList.Achis[i].AccuracyStage = stage
				currentList.Achis[i].AccuracyStageNext = stageNext
			} else if stageName == "prec" {
				currentList.Achis[i].PrecisionStage = stage
				currentList.Achis[i].PrecisionStageNext = stageNext
			} else if stageName == "read" {
				currentList.Achis[i].ReadingStage = stage
				currentList.Achis[i].ReadingStageNext = stageNext
			} else if stageName == "speed" {
				currentList.Achis[i].SpeedStage = stage
				currentList.Achis[i].SpeedStageNext = stageNext
			} else if stageName == "stam" {
				currentList.Achis[i].StaminaStage = stage
				currentList.Achis[i].StaminaStageNext = stageNext
			}
		}

	}

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/achi.json", finalList, 0644)

}

func setTutorialDone(id int) {
	currentList := GetAchiJSON()

	for i := 0; i < len(currentList.Achis); i++ {
		if currentList.Achis[i].Id == id {
			currentList.Achis[i].Stage = 1
			currentList.Achis[i].AccuracyStage = 1
			currentList.Achis[i].AccuracyStageNext = "99% acc on any map."
			currentList.Achis[i].PrecisionStage = 1
			currentList.Achis[i].ReadingStage = 1
			currentList.Achis[i].SpeedStage = 1
			currentList.Achis[i].StaminaStage = 1
		}
	}

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/achi.json", finalList, 0644)
}

func GetAchi(id int) Achi {
	allAchis := GetAchiJSON()
	for i := 0; i < len(allAchis.Achis); i++ {
		if allAchis.Achis[i].Id == id {
			return allAchis.Achis[i]
		}
	}
	return Achi{}
}

func NewAchi(newUser Achi) {
	currentList := GetAchiJSON()

	currentList.Achis = append(currentList.Achis, Achi{
		Id:    newUser.Id,
		Stage: 0,
	})

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/achi.json", finalList, 0644)
}

func GetAchiJSON() Achis {
	jsonFile, err := os.Open("web/data/achi.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var achi Achis

	json.Unmarshal(byteValue, &achi)

	return achi
}
