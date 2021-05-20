package Infrastructure

import (
	"TwitchChatBot/Infrastructure/mock_Infrastructure"
	"log"
	"testing"
	"time"
)

type rateLimiterTestHarness struct {
	timeObject   *mock_Infrastructure.FakeITime
	rateLimit    int
	rateDuration time.Duration
	patient      IRateLimiter
}

func setupRateLimiterTestHarness() *rateLimiterTestHarness {
	testHarness := new(rateLimiterTestHarness)
	testHarness.timeObject = mock_Infrastructure.NewFakeITime()
	testHarness.rateLimit = 3
	testHarness.rateDuration = time.Minute
	testHarness.patient = NewRateLimiter(testHarness.rateLimit, testHarness.rateDuration, testHarness.timeObject)
	return testHarness
}

func Test_WillPerformInteractionImmediatelyIfInterationAvailable(test *testing.T) {
	testHarness := setupRateLimiterTestHarness()

	complete := false
	go testHarness.patient.PerformInteraction(
		func() {
			complete = true
		},
	)
	time.Sleep(time.Millisecond)

	if !complete {
		test.Errorf("Interaction was not performed")
	}
}

func Test_WillNotPerformInteractionImmediatelyIfInteractionNotAvailable(test *testing.T) {
	testHarness := setupRateLimiterTestHarness()
	for i := 0; i < testHarness.rateLimit; i++ {
		testHarness.patient.PerformInteraction(
			func() {},
		)
	}

	complete := false
	go testHarness.patient.PerformInteraction(
		func() {
			complete = true
		},
	)
	time.Sleep(time.Millisecond)

	if complete {
		test.Errorf("Interaction was performed unexpectedly")
	}
}

func Test_WillNotPerformInteractionBeforeRequiredTimeElapses(test *testing.T) {
	testHarness := setupRateLimiterTestHarness()
	for i := 0; i < testHarness.rateLimit; i++ {
		testHarness.patient.PerformInteraction(
			func() {
				time.Sleep(time.Millisecond)
			},
		)
		testHarness.timeObject.AdvanceTime(time.Second)
	}

	complete := false
	go testHarness.patient.PerformInteraction(
		func() {
			complete = true
			time.Sleep(time.Millisecond)
		},
	)
	log.Println("About to advance time")
	testHarness.timeObject.AdvanceTime(testHarness.rateDuration - time.Second*10)
	time.Sleep(time.Millisecond)

	if complete {
		test.Errorf("Interaction was performed unexpectedly")
	}
}

func Test_WillPerformInteractionAfterRequiredTimeElapses(test *testing.T) {
	testHarness := setupRateLimiterTestHarness()
	for i := 0; i < testHarness.rateLimit; i++ {
		testHarness.patient.PerformInteraction(
			func() {
				time.Sleep(time.Millisecond)
			},
		)
		testHarness.timeObject.AdvanceTime(time.Second)
	}

	complete := false
	go testHarness.patient.PerformInteraction(
		func() {
			complete = true
			time.Sleep(time.Millisecond)
		},
	)
	testHarness.timeObject.AdvanceTime(testHarness.rateDuration + time.Second)
	time.Sleep(time.Millisecond)

	if !complete {
		test.Errorf("Interaction was not performed")
	}
}

func Test_WillPerformAllQueuedInteractionsOnceEnoughTimePasses(test *testing.T) {
	log.Println("Starting desired test")
	testHarness := setupRateLimiterTestHarness()
	for i := 0; i < testHarness.rateLimit; i++ {
		testHarness.patient.PerformInteraction(
			func() {
				time.Sleep(time.Millisecond)
			},
		)
		testHarness.timeObject.AdvanceTime(time.Second)
	}

	completeCount := 0
	for i := 0; i < 3; i++ {
		go testHarness.patient.PerformInteraction(
			func() {
				completeCount++
				time.Sleep(time.Millisecond)
			},
		)
	}
	log.Println("about to advance time")
	testHarness.timeObject.AdvanceTime(testHarness.rateDuration + time.Second)
	time.Sleep(time.Millisecond)

	if completeCount != 3 {
		test.Errorf("Only %d queued interactions were performed", completeCount)
	}
}

func Test_WillNotPerformInteractionAfterShutdown(test *testing.T) {
	testHarness := setupRateLimiterTestHarness()

	complete := false
	go testHarness.patient.PerformInteraction(
		func() {
			complete = true
			time.Sleep(time.Millisecond)
		},
	)
	testHarness.patient.ShutDown()

	testHarness.timeObject.AdvanceTime(testHarness.rateDuration + time.Second)
	time.Sleep(time.Millisecond)

	if complete {
		test.Errorf("Interaction was performed unexpectedly")
	}
}

func Test_WillNotPerformQueuedInteractionAfterShutdown(test *testing.T) {
	testHarness := setupRateLimiterTestHarness()
	for i := 0; i < testHarness.rateLimit; i++ {
		testHarness.patient.PerformInteraction(
			func() {
				time.Sleep(time.Millisecond)
			},
		)
		testHarness.timeObject.AdvanceTime(time.Second)
	}

	complete := false
	go testHarness.patient.PerformInteraction(
		func() {
			complete = true
			time.Sleep(time.Millisecond)
		},
	)
	testHarness.patient.ShutDown()

	testHarness.timeObject.AdvanceTime(testHarness.rateDuration + time.Second)
	time.Sleep(time.Millisecond)

	if complete {
		test.Errorf("Interaction was performed unexpectedly")
	}
}
