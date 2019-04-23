# RxJS TestScheduler - A Typescript and Jest Based Example

## Overview

RxJs is the new cool tool in the 


 has been a crucial part of our tech stack for over a year now. We use this library to fetch, and manipulate streams of data. If you want to learn more about rxjs in general checkout [some](https://vendasta-blog.appspot.com/blog/BL-TRRN3PBW/) [of](https://vendasta-blog.appspot.com/blog/BL-KJ8SCZ5C/) [these](https://vendasta-blog.appspot.com/blog/BL-VKC8SC2T/) [blogs](https://vendasta-blog.appspot.com/blog/BL-LBQ7XWWH/)! This blog will not be about how to use rxjs, but rather **how to unit test typescript code that uses rxjs observables**. The goal of this blog is to provide you with the base knowledge you need to start testing functionality that depends on observables.

### Observable Review

During runtime, an [observable](http://reactivex.io/documentation/observable.html) is a object that emits entities as a stream of data. Observables also contain other information like;

- stream completion, and;
- error information.

There are two different observables. A cold observable is one that doesn't emit any values until a subscription is performed onto it. While a hot observable does. A good example of a cold observable is one that comes from the http service. An example of a hot observable would be from a rxjs subject.

### When Not To Use It

The rxjs test scheduler isn't designed to be used to test the return observable of a unit without the input being a observable controlled by it. Let's say you just made a shiny new function that returns a new observable:

``` typescript
function someThing(someString: string) {
  // Do some operation on someString
  return Observable.of(someString)
}
```

Using jasmine as the testing framework, you could test this function with:

``` typescript
it('returns hello world for hello world string input', () => {
  let result = someThing('hello world')
    .subscribe((val) => expect(val).toBe('hello world'));
})
```

Note that asynchronous assertions are not ideal in general. Unless you know an event will be emitted from the observable, don't assert in a subscribe.

## Testscheduler Basics

The rxjs test scheduler is designed to compare the characteristics of two observables. A simple test would have the following operations:

1; Create a new observable that is in the control of the test scheduler.
2; Use that observable as the input into your unit.
3; Compare the output observable with what you are expecting it to be.

### Setup

Setup is actually really easy! All you need to do is import test scheduler and then create a new instance of it.

``` typescript
import {TestScheduler} from 'rxjs';

describe('SomeThing', () => {
  ...
  let testScheduler = new TestScheduler((a,b) => expect(a).toEqual(b))
})
```

Note that the required parameter of the `TestScheduler` constructor is a callback function to do the assertions on. Since we use jasmine, our callback will look like `(a,b) => expect(a).toEqual(b)`. The reason why we wouldn't want to use jasmine's `.toBe()` method function, is that it does a strict comparison. If the observables exist in different locations in memory your assertion would fail. This is something that has burned me.

### Creating An Observable

To create an observable with the test scheduler we need to define two things:

1; The emitted events, errors, and stream completions on a linear time scale.
2; The emitter event values

An example of this is:

``` typescript
  let input = testScheduler.createColdObservable('--a--#--|', {a: SomeObject/Primative})
```

You can also create a hot observable with the same syntax. Whatever you choose to create should reflect the actual input into your unit.

### Marbles

You were probably thinking, "what the heck is `'--a--#--|'`"? That is actually a marble diagram. Since we need to do comparisons in discrete time, we have to specify our specify our events in terms of marbles.

Lets break this down. First off, each hyphen represents 10 "units" of time. The "a" is actually an event, and the line signifies the end of a stream (`Observable.empty()`). The pound sign is an `Observable.throw()`. For the official guide click [here](https://github.com/ReactiveX/rxjs/blob/master/doc/writing-marble-tests.md).

### Assertions

The test scheduler has an assertion mechanism that wraps a projects testing framework. As mentioned before, it's characteristics are defined by the callback function passed into the parameter of the constructor.

An example of this is:

``` typescript
  testScheduler.expectObservable(testComponent.someObservable).toBe('--a--|', {a: someObject});
  testScheduler.flush();
```

Note that `TestScheduler.flush()` executes all discrete observable streams with the testScheduler and must be called after a `TestScheduler.expectObservable` assertion.

### Practical Tests

For real world examples of how we are currently using the test scheduler, see:

- [marketplace public store](https://github.com/vendasta/marketplace-public-store/blob/master/src/app/package-page/interest-in-package-dialog/interest-in-package-dialog.component.spec.ts)
- [partner center client](https://github.com/vendasta/partner-center-client/blob/18752024dc205a5437f59d219a3abfe519c956ff/src/app/sales-orders/sales-order-details.service.spec.ts)
- [sales tool](https://github.com/vendasta/ST/search?utf8=%E2%9C%93&q=%22Test+Scheduler%22&type=)

### Example

In the following sections I will be going over some different tests we might want to perform on a unit that contains observable functionality. I will be using an example in production which you can find [here](https://github.com/vendasta/vendor-center-client/blob/master/src/angular/partner/details/partner-details.component.spec.ts). The unit that we will be examining is:

``` typescript
public initializePartnerStream(): void {
  this.partner$ = this.partnerService
    .get(this.partnerId)
    .catch(() => {
      this.apiFailed = true;
      this.alertService.errorSnack('Oops! We were not able to load any data');
      return Observable.of(null);
});
```

#### Assignments

We want to test that on a successful api call from the `partnerService`, an observable containing a partner is returned.

```typescript
it('should return an object representing a partner when the api call succeeds', () => {
  partnerServiceSpy.get.and.callFake(() => testScheduler.createColdObservable('--a--|', {a: testPartner}));
  testComponent = new PartnerDetailsComponent(dialogSpy, partnerServiceSpy, alertServiceSpy);
  testComponent.initializePartnerStream();
  testScheduler.expectObservable(testComponent.partner$).toBe('--a--|', {a: testPartner});
  testScheduler.flush();
});
```

#### Catching Errors

Now we should check to see that the `errorSnack` Was thrown.

```typescript
it('should call the error snack method in the alert service', () => {
      partnerServiceSpy.get.and.callFake(() => {
        return testScheduler.createHotObservable('#');
      });
      testComponent = new PartnerDetailsComponent(dialogSpy, partnerServiceSpy, alertServiceSpy);
      testComponent.initializePartnerStream();
      testScheduler.flush();
      testComponent.partner$.subscribe((result) => {
        expect(alertServiceSpy.errorSnack).toHaveBeenCalled();
      });
    });
});
```
