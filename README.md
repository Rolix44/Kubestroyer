<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a name="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
<div align="center">

  [![Contributors][contributors-shield]][contributors-url]
  [![Forks][forks-shield]][forks-url]
  [![Stargazers][stars-shield]][stars-url]
  [![Issues][issues-shield]][issues-url]
  [![MIT License][license-shield]][license-url]

</div>



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/Rolix44/Kubestroyer">
    <img src=".github/logo.png" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">Kubestroyer</h3>

  <p align="center">
    Kubestroyer aims to exploit Kubernetes clusters misconfigurations and be the swiss army knife of your Kubernetes pentests
    <br />
    <a href="https://github.com/Rolix44/Kubestroyer"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/Rolix44/Kubestroyer/issues">Report Bug</a>
    ·
    <a href="https://github.com/Rolix44/Kubestroyer/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

[![Product Name Screen Shot][product-screenshot]]()

**Kubestroyer** is a Golang exploitation tool that aims to take advantage of Kubernetes clusters misconfigurations.

The tool is scanning known Kubernetes ports that can be exposed as well as exploiting them.


<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With
<br>

[![Golang][Golang]][Golang-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running, follow these simple example steps.

### Prerequisites

* Go 1.19
  ```sh
  wget https://go.dev/dl/go1.19.4.linux-amd64.tar.gz
  tar -C /usr/local -xzf go1.19.4.linux-amd64.tar.gz
  ```

### Installation
Use [prebuilt binary](https://github.com/Rolix44/Kubestroyer/releases) or build from source:

1. Clone the repo
   ```sh
   git clone https://github.com/Rolix44/Kubestroyer.git
   ```
2. build the binary
   ```sh
   go build -o Kubestroyer cmd/kubestroyer/main.go 
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage
[![Product usage][usage-gif]]()

|  Parameter | Description | Mand/opt  | Example  | 
|---|---|---|---|
|  -t / --target | Target (IP, domain or file)  |  Mandatory | -t localhost,127.0.0.1 /  -t ./domain.txt  |
| --node-scan  | Enable node port scanning (port 30000 to 32767)  |  Optionnal | -t localhost --node-scan  |
| --anon-rce  | RCE using Kubelet API anonymous auth  | Optionnal  | -t localhost --anon-rce  |
|  -x  |  Command to execute when using RCE (display service account token by default) | Optionnal  | -t localhost --anon-rce -x "ls -al"  |


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [x] Split main into different packages
- [x] Target selection
    - [x] List of multiple targets
    - [x] input file as target parameter
- [ ] Description for each detected open port
- [ ] Choose the pod for anon RCE 

See the [open issues](https://github.com/Rolix44/Kubestroyer/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Rolix - [@Rolix_cy](https://twitter.com/Rolix_cy) - rolixcy@protonmail.com

Project Link: [https://github.com/Rolix44/Kubestroyer](https://github.com/Rolix44/Kubestroyer)

<p align="right">(<a href="#readme-top">back to top</a>)</p>




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/Rolix44/Kubestroyer.svg?style=for-the-badge
[contributors-url]: https://github.com/Rolix44/Kubestroyer/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/Rolix44/Kubestroyer.svg?style=for-the-badge
[forks-url]: https://github.com/Rolix44/Kubestroyer/network/members
[stars-shield]: https://img.shields.io/github/stars/Rolix44/Kubestroyer.svg?style=for-the-badge
[stars-url]: https://github.com/Rolix44/Kubestroyer/stargazers
[issues-shield]: https://img.shields.io/github/issues/Rolix44/Kubestroyer.svg?style=for-the-badge
[issues-url]: https://github.com/Rolix44/Kubestroyer/issues
[license-shield]: https://img.shields.io/github/license/Rolix44/Kubestroyer.svg?style=for-the-badge
[license-url]: https://github.com/Rolix44/Kubestroyer/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[Golang]: https://img.shields.io/badge/Golang-0769AD?style=for-the-badge&logo=Go&logoColor=white
[Golang-url]: https://go.dev/
[product-screenshot]: .github/project-name.png
[release]: https://github.com/Rolix44/Kubestroyer/releases/new
[usage-gif]: .github/kubestroyer-usage.gif