# A Prometheus metrics exporter for Pika

A simple plugin to export stats from Pika to Prometheus. This plugin acts as a
proxy to the Pika JSONRPC method `session.stats` and re-exports it in a Prometheus
compatible format.
