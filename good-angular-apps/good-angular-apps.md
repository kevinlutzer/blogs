# How to Take your Angular Apps to the Next Level


## Consistently Use Loading, Error, and "No Data" View for Each Page. 

How does your page component react when there is no data to show? What happens if the api you are fetching data from fails?

## Use The `noscript` Tag

The `noscript` tag is an html5 element whose child elements are parsed for users who have javascript disabled on their browser. The content but between the opening and closing tags can be any combination of html elements. We can use this tag to define custom messages to encourage our users to enable javascript. Note that currently for all of our angular apps, our users would just see a blank white page. Here is a simple example:

``` html
...
<body>
  <noscript>It appears that you have javascript disabled. Consider enabling it to have access to this website's features.</noscript>
</body>
...
```
In the browser this would look like: 
<INSERT IMAGE HERE>


## Lazy Load Modules




## Add a Splash Screen

When the index and bundle files for an angular app are loaded by a user, there is a period of time where the user will not see any app content. Why? Well that period of time corresponds to angular boostrapping your app on the clients computer. Depending on the size of your main bundle, the functionality in the app's entry point file `main.ts`, and a lot of other factors, your app could take almost a second to bootstrap. To improve the user experience we could provide some sort of splash screen as a temporary substitute.

Say that your `body` tag in the app's `index.html` file looks something like: 

``` html
  <body>
    <app-root> </app-root>
  </body>
```

What we could do is add an image tag and some custom styling inbetween the opening and closing `app-root` tags like so: 

``` html
  <body>
    <app-root> 
      <img src="some-gcs-bucket.com/some/variation/of/the/vendasta/logo" styles="display:flex; justify-contents: center; align-items; center"><img>
    </app-root>
  </body>
```

This produces the following view. INSERT IMAGE HERE
An example usage of this idea is the [mission control client](linkToGithubUrl) application. 

