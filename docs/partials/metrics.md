jenkins_job_buildable{name, path, class}
: 1 if the job is buildable, 0 otherwise

jenkins_job_color{name, path, class}
: Color code of the jenkins job

jenkins_job_disabled{name, path, class}
: 1 if the job is disabled, 0 otherwise

jenkins_job_duration{name, path, class}
: Duration of last build in ms

jenkins_job_end_time{name, path, class}
: Start time of last build as unix timestamp

jenkins_job_last_build{name, path, class}
: Builder number for last build

jenkins_job_last_completed_build{name, path, class}
: Builder number for last completed build

jenkins_job_last_failed_build{name, path, class}
: Builder number for last failed build

jenkins_job_last_stable_build{name, path, class}
: Builder number for last stable build

jenkins_job_last_successful_build{name, path, class}
: Builder number for last successful build

jenkins_job_last_unstable_build{name, path, class}
: Builder number for last unstable build

jenkins_job_last_unsuccessful_build{name, path, class}
: Builder number for last unsuccessful build

jenkins_job_next_build_number{name, path, class}
: Next build number for the job

jenkins_job_start_time{name, path, class}
: Start time of last build as unix timestamp

jenkins_request_duration_seconds{collector}
: Histogram of latencies for requests to the api per collector

jenkins_request_failures_total{collector}
: Total number of failed requests to the api per collector
