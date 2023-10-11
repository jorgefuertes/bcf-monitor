package common

import "time"

type MonitorizableBase struct {
	fails                int
	failedAt             time.Time
	notifiedAt           time.Time
	TimeoutSeconds       int
	EverySeconds         int
	FailAfterRetries     int
	RememberAfterMinutes int
}

func (s *MonitorizableBase) AddFail() {
	if s.fails == 0 {
		s.failedAt = time.Now()
	}
	s.fails++
}

func (s *MonitorizableBase) IsDown() bool {
	return s.Fails() >= s.FailAfterRetries
}

func (s *MonitorizableBase) IsUp() bool { return !s.IsDown() }

func (s *MonitorizableBase) Fails() int { return s.fails }
func (s *MonitorizableBase) FailedAt() time.Time { return s.failedAt }

func (s *MonitorizableBase) SetNotifiedNow() { s.notifiedAt = time.Now() }
func (s *MonitorizableBase) ClearNotified() { s.notifiedAt = time.Time{} }
func (s *MonitorizableBase) IsNotified() bool { return !s.notifiedAt.IsZero() }

func (s *MonitorizableBase) GetNotifiedAt() time.Time { return s.notifiedAt }

func (s *MonitorizableBase) SinceLastNotification() time.Duration { return time.Until(s.GetNotifiedAt()) }

func (s *MonitorizableBase) Reset() {
	s.fails = 0
	s.ClearNotified()
}

func (s *MonitorizableBase) Every() time.Duration {
	return time.Duration(s.EverySeconds) * time.Second
}

func (s *MonitorizableBase) Timeout() time.Duration {
	return time.Duration(s.TimeoutSeconds) * time.Second
}

func (s *MonitorizableBase) RememberAfter() time.Duration {
	return time.Duration(s.RememberAfterMinutes) * time.Minute
}

func (s *MonitorizableBase) CanBeNotified() bool {
	return s.SinceLastNotification() > s.RememberAfter()
}
