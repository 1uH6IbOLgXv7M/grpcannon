// Package pause implements a lightweight pause/resume gate for worker
// goroutines. Workers call Wait before each unit of work; the gate blocks
// them while the controller is paused and releases them atomically on Resume.
//
// Typical usage:
//
//	ctrl := pause.New()
//	go func() {
//		for {
//			ctrl.Wait() // blocks if paused
//			doWork()
//		}
//	}()
//
//	ctrl.Pause()
//	// ... reconfigure ...
//	ctrl.Resume()
package pause
