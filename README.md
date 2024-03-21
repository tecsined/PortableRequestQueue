# Request Processing with Tracking and Retry

This program facilitates processing a list of HTTP requests with tracking and retry mechanisms. It ensures reliable execution and avoids repetitive work on successful requests.

## Prerequisites

* Go programming language installed (https://go.dev/doc/install)

## Getting Started

1. **Delete Tracking File (Optional):**
   If you have a previous run's `completed_requests.json` file, it's recommended to delete it to ensure a clean start:

   ```bash
   rm completed_requests.json

