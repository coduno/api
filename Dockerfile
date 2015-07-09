# See:
#
#  Building Custom Runtimes: Base Images
#  https://cloud.google.com/appengine/docs/managed-vms/custom-runtimes#base_images
#
#  Developing and Deploying Managed VMs: Dockerfiles
#  https://cloud.google.com/appengine/docs/managed-vms/sdk#dockerfiles

FROM gcr.io/google_appengine/golang

COPY . /go/src/app

RUN go-wrapper download
RUN go-wrapper install -tags appenginevm
