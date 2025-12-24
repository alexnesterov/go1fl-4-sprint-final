package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	s := strings.Split(data, ",")
	if len(s) != 3 {
		return 0, "", 0, errors.New("error parsing string")
	}

	steps, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("number of steps should be greater than 0")
	}

	activity := s[1]

	duration, err := time.ParseDuration(s[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration should be greater than 0")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distance := float64(steps) * stepLength / mInKm

	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)
	speed := distance / duration.Hours()

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	var resultBuilder strings.Builder

	steps, activity, durarion, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	fmt.Fprintf(&resultBuilder, "Тип тренировки: %s\n", activity)
	fmt.Fprintf(&resultBuilder, "Длительность: %.2f ч.\n", durarion.Hours())

	distance := distance(steps, height)
	fmt.Fprintf(&resultBuilder, "Дистанция: %.2f км.\n", distance)

	speed := meanSpeed(steps, height, durarion)
	fmt.Fprintf(&resultBuilder, "Скорость: %.2f км/ч\n", speed)

	var calories float64

	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, durarion)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, durarion)
	default:
		return "", errors.New("unknown type of training")
	}

	if err != nil {
		return "", err
	}

	fmt.Fprintf(&resultBuilder, "Сожгли калорий: %.2f\n", calories)

	return resultBuilder.String(), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, errors.New("weight should be greater than 0")
	}
	if steps <= 0 {
		return 0, errors.New("number of steps should be greater than 0")
	}
	if duration <= 0 {
		return 0, errors.New("duration should be greater than 0")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := weight * meanSpeed * durationInMinutes / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, errors.New("weight should be greater than 0")
	}
	if height <= 0 {
		return 0, errors.New("height should be greater than 0")
	}
	if steps <= 0 {
		return 0, errors.New("number of steps should be greater than 0")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := weight * meanSpeed * durationInMinutes / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
