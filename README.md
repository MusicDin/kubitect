<div align=center>
<h1><a href="https://kubitect.io">Kubitect</a></h1>
<img src="docs/theme/assets/images/favicon.svg" width=200></img>

</br>

# What is Kubitect?

Kubitect is an open source project that aims to simplify the **deployment** and subsequent **management of Kubernetes clusters**.
It provides a CLI tool written in Golang that lets you **set up**, **upgrade**, **scale**, and **destroy** Kubernetes clusters.
Under the hood, it uses [Terraform](https://www.terraform.io/) along with [terraform-libvirt-provider](https://github.com/dmacvicar/terraform-provider-libvirt)
to deploy virtual machines on target hosts running libvirt.
Kubernetes is configured on the deployed virtual machines using [Kubespray](https://kubespray.io), the popular open source project.

</br>

![Go Test Badge](https://github.com/marmiha/kubitect/actions/workflows/go-test.yml/badge.svg)
![Release Badge](https://github.com/marmiha/kubitect/actions/workflows/release-cli-binaries.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/MusicDin/kubitect.svg)](https://pkg.go.dev/github.com/MusicDin/kubitect)

### Documentation

Kubitect documentation is accessible at **[:book: kubitect.io](https://kubitect.io/getting-started/installation/)**.

![Docs Workflow Badge](https://github.com/marmiha/kubitect/actions/workflows/docs.yml/badge.svg)
[![MkDocs](https://img.shields.io/badge/docs-kubitect.io-blue)](https://kubitect.io)

</br>

### Releases

All Kubitect releases are available on the [release page](https://github.com/MusicDin/kubitect/releases).
</br>
It is recommended to use official releases, as unreleased versions from the master branch may be unstable.

</br>

### Authors

[Din Mušić](https://github.com/MusicDin) and [all contributors](https://github.com/MusicDin/kubitect/graphs/contributors).

</br>

### License

[Apache License 2.0](./LICENSE)

</div>
