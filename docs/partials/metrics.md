jenkins_job_buildable{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: 1 if the sjob is buildable, 0 otherwise

jenkins_job_color{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Color code of the jenkins job

jenkins_job_disabled{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: 1 if the job is disabled, 0 otherwise

jenkins_job_last_build{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Builder number for last build

jenkins_job_last_completed_build{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Builder number for last completed build

jenkins_job_last_failed_build{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Builder number for last failed build

jenkins_job_last_stable_build{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Builder number for last stable build

jenkins_job_last_successful_build{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Builder number for last successful build

jenkins_job_last_unstable_build{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Builder number for last unstable build

jenkins_job_last_unsuccessful_build{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Builder number for last unsuccessful build

jenkins_job_next_build_number{<prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>, <prometheus.ConstrainedLabel Value>}
: Next build number for the job

jenkins_request_duration_seconds{collector}
: Histogram of latencies for requests to the api per collector

jenkins_request_failures_total{collector}
: Total number of failed requests to the api per collector
