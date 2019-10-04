

# Intro

Vendasta publishes shared Angular 2 components to our npm repository to create a single source of truth. How do we test and develop these packages locally?

## Nomenclature

Within this guide I will use the term 'test package' to describe a package that is under development. The term 'test project' will be used to represent the hypothetical project we would test our test package with. In a practical setting, a test package may be something like `@vendasta/package-details`, and a test project would be something like `marketplace-public-store`.

This guide discusses a test package which has the following structure:

<INSERT DIAGRAM>

The following contents may not fully apply to packages like `@vendasta/uikit`.

This guide assumes that you have npm v5 installed on your computer. If you don't, install it by running `npm install npm@latest -g`.

## Release Candidates

One way to develop a test package is to use release candidates. A release candidate is a special type of version that npm provides through its versioning semantics. In a general sense, a release candidate is meant to be purely a test version of any package. The code within a release candidate is not guaranteed to be production ready. Publishing and then installing release candidates for any new changes to the test package is a taxing process as there is no automatic change detection. 

Currently, the exhausting process of using release candidates is as follows:

1. Set a new release candidate version in the package.json file for the test package. This version will look something like `X.X.X-rc.XX`.

2. Run the angular compiler (`ngc`) in the test package's root directory. This will generate the metadata, declaration, and javascript files for your test package's modules. If you don't have the angular compiler with it's dependencies installed globally, you can source it from any of our front-end projects in `node_modules/.bin/`.

3. Publish the release candidate. To do this run `npm publish` in the package's root directory.

4. Install the newest release candidate into the project you are testing your package with. Usually this involves running `npm upgrade @npm-user/package-name`. Most likely your running webpack server will pick up these changes. If not, stop the webpack server instance and then make sure that the correct release candidate of the package was installed; start the server again.

After you have completed these steps four or five times you will start to wish that there was a better way. Well there is!

# Npm Link

Instead of publishing release candidates of a test package, we can just source the local test package code using npm link. The following diagram demonstrates the relationship npm link has with npm publish.

ADD DIAGRAM

When you run `npm link` in your test package's root directory, it creates a symbolic link from your global node modules to the test package's source code. You can source this test package through the symbolic link by running `npm link @npm-user/package-name` in the root directory of the test project you want to try your test package in. Npm link will automatically copy the code from your packages current project to `node_modules/@npm-user/package-name` of the test project. Now you will have change detection! 

The steps for setting up npm link are very simple:

1. In the test package's root directory run `npm link`
2. In the root directory of the test project run `npm link @npm-user/package-name`
3. Start your webpack server instance if you haven't already

# Config

You might need to do a little extra configuration to get npm link working with your test project.

One of the issues you may run into is that the webpack server instance might source your test package's dependencies from the project it is in, rather then the project you are testing it with. To fix this add the following json property and value to the tsconfig of your test project.

``` json
{
    ...
    baseUrl: '<source dir>`
    paths: {
        '@user-name/package-name' : [
            '<path to node_modules from baseUrl>/node_modules/@user-name/package-name`
        ]
    }
    ...
}
```

Also delete the `paths` object in the package.json file and the node modules in the test package.

You should have everything you need to get up and running with npm link!

