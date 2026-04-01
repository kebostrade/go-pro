{{/*
Expand the name of the chart.
*/}}
{{- define "gopro-chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "gopro-chart.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gopro-chart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels for gopro-chart.
*/}}
{{- define "gopro-chart.labels" -}}
app.kubernetes.io/name: {{ include "gopro-chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/part-of: {{ .Chart.Name }}
helm.sh/chart: {{ include "gopro-chart.chart" . }}
{{- end }}

{{/*
Selector labels for gopro-chart.
*/}}
{{- define "gopro-chart.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gopro-chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
