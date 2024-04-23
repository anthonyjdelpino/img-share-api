
 <div align="center">
 <a href="https://github.com/anthonyjdelpino/img-share-api">
    <img src="logo.svg" alt="Logo" width="80" height="80">
  </a>

# img-share-api
</div>

![deploystatus](https://github.com/anthonyjdelpino/img-share-api/actions/workflows/main.yml/badge.svg?branch=prod)


img-share-api is a RESTful image sharing API written in GoLang, hosted on Google Cloud Functions, and deployed with GitHub Actions :)

### Built With
* ![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)
* ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
* ![Google Cloud](https://img.shields.io/badge/GoogleCloud-%234285F4.svg?style=for-the-badge&logo=google-cloud&logoColor=white)

&nbsp;

## Usage

You can use a browser or curl to invoke the API at the current endpoint: 

[https://us-central1-img-share-api-project.cloudfunctions.net/img-share-api-func](https://us-central1-img-share-api-project.cloudfunctions.net/img-share-api-func)

The following are the API endpoints:

`/`

`/images`

`/images/<id>`

&nbsp;

To get all image listings, send a GET request to the `/images` endpoint:

```bash
curl https://us-central1-img-share-api-project.cloudfunctions.net/img-share-api-func/images
```
To get a specific image listing, send a GET request to the `images/<id>` endpoint, where `<id>` is the ID of the image:

```bash
curl https://us-central1-img-share-api-project.cloudfunctions.net/img-share-api-func/images/IMAGE-ID
```

To upload an image, send a POST request to the `/images` endpoint:
```bash
curl -X POST "https://us-central1-img-share-api-project.cloudfunctions.net/img-share-api-func/images" -F file=@/PATH/TO/YOUR/IMAGE
```

To delete an image, send a DELETE request to the `images/<id>` endpoint, where `<id>` is the ID of the image:

```bash
curl -X DELETE "https://us-central1-img-share-api-project.cloudfunctions.net/img-share-api-func/images/IMAGE-ID"
```