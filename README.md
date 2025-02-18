<a name="readme-top"></a>


<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/alexanderbkl/vidre-back">
    <img src="assets/imgs/logo.jpg" alt="Logo" width="100" height="100">
  </a>

<h3 align="center">VidreBany (Backend)</h3>

  <p align="center">
    Factory management system
    <br />
    <br />
    ·
    <a href="https://github.com/alexanderbkl/vidre-back/issues">Report Bug</a>
    ·
    <a href="https://github.com/alexanderbkl/vidre-back/issues">Request Feature</a>
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
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project
**[vidrebany.com](https://vidrebany.com)**
<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

[![Go][Go]][Go-url]
[![Postgres][Postgres]][Postgres-url]
[![Docker][Docker]][Docker-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### ✔️ Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* go (version 1.15 or over)
* docker desktop
* make (install by choco)

### 🚀 Installation

1. Clone the repo
   ```sh
   git clone https://github.com/alexanderbkl/vidre-back.git
   ```
2. create <b>.env</b> file from <b>.env.example</b>
3. run docker
   ```sh
   make production
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## 💡 Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## 🤠 Roadmap

- [ ] Feature 1
- [ ] Feature 2
- [ ] Feature 3
    - [ ] Nested Feature

See the [open issues](https://github.com/alexanderbkl/vidre-back/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## 🤝 Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again! ⭐⭐⭐

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feat/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feat/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<p>REMOVE docker backend-dev container logs:</p>

```sh
sudo -i &&
echo "" > $(docker inspect --format='{{.LogPath}}' backend-dev)
```
<!-- LICENSE -->
## 📝 License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>





<!-- ACKNOWLEDGMENTS -->
## 📒 Acknowledgments

* commented todo list like following
```
// TO-DO *****
```


<p align="right">(<a href="#readme-top">back to top</a>)</p>


## 🚧 Under Building 🚧

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[product-screenshot]: images/screenshot.png
[Go]: https://img.shields.io/badge/go-20232A?style=for-the-badge&logo=go
[Go-url]: https://go.dev/
[Postgres]: https://img.shields.io/badge/Postgres-20232A?style=for-the-badge&logo=postgresql&logoColor=61DAFB
[Postgres-url]: https://www.postgresql.org/
[Docker]: https://img.shields.io/badge/docker-20232A?style=for-the-badge&logo=docker
[Docker-url]: https://docker.com/
