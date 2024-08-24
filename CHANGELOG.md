# Changelog for 1.1.0

The following sections list the changes for 1.1.0.

## Summary

 * Chg #51: Read secrets form files
 * Chg #51: Integrate standard web config
 * Enh #51: Integrate option pprof profiling

## Details

 * Change #51: Read secrets form files

   We have added proper support to load secrets like the password from files or
   from base64-encoded strings. Just provide the flags or environment variables for
   token or private key with a DSN formatted string like `file://path/to/file` or
   `base64://Zm9vYmFy`.

   https://github.com/promhippie/jenkins_exporter/pull/51

 * Change #51: Integrate standard web config

   We integrated the new web config from the Prometheus toolkit which provides a
   configuration for TLS support and also some basic builtin authentication. For
   the detailed configuration you can check out the documentation.

   https://github.com/promhippie/jenkins_exporter/pull/51

 * Enhancement #51: Integrate option pprof profiling

   We have added an option to enable a pprof endpoint for proper profiling support
   with the help of tools like Parca. The endpoint `/debug/pprof` can now
   optionally be enabled to get the profiling details for catching potential memory
   leaks.

   https://github.com/promhippie/jenkins_exporter/pull/51


# Changelog for 1.0.0

The following sections list the changes for 1.0.0.

## Summary

 * Chg #6: Initial release of basic version

## Details

 * Change #6: Initial release of basic version

   Just prepared an initial basic version which could be released to the public.

   https://github.com/promhippie/jenkins_exporter/issues/6


