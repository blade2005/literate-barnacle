# 20240620. Frontend Implementation Language

Date: 2024-06-24

## Status

Proposed

## Context

Deciding how to implement the frontend.

### Approaches considered

#### WebAssembly(WASM)

In general WASM is [supported](https://developer.mozilla.org/en-US/docs/WebAssembly#browser_compatibility) mostly by most browsers. There will be some "gotchas" in writing the code.

From the [Mozila MDN](https://developer.mozilla.org/en-US/docs/WebAssembly)


> WebAssembly is a type of code that can be run in modern web browsers — it is a low-level assembly-like language with a compact binary format that runs with near-native performance and provides languages such as C/C++, C# and Rust with a compilation target so that they can run on the web. It is also designed to run alongside JavaScript, allowing both to work together.
>
> WebAssembly has huge implications for the web platform — it provides a way to run code written in multiple languages on the web at near-native speed, with client apps running on the web that previously couldn't have done so.
>
> WebAssembly is designed to complement and run alongside JavaScript — using the WebAssembly JavaScript APIs, you can load WebAssembly modules into a JavaScript app and share functionality between the two. This allows you to take advantage of WebAssembly's performance and power and JavaScript's expressiveness and flexibility in the same app, even if you don't know how to write WebAssembly code.
>
> And what's even better is that it is being developed as a web standard via the W3C WebAssembly Working Group and Community Group with active participation from all major browser vendors.

WASM code still needs to be fetched and executed by Javascript, so we'll still have a tiny portion of js code to have in the project, but this would be likely write it once and never touch it again.

Go has a section talking about [WebAssembly](https://go.dev/wiki/WebAssembly) where I got a number of the project mentioned below from.

##### Vecty

[Vecty](https://github.com/hexops/vecty) is written in Go. It is still considered an experimental work-in-progress. They provide a list of items not yet [ready](https://github.com/hexops/vecty?tab=readme-ov-file#current-status)

The biggest issue I would forsee if the lack of a ready to use component library, this would mean needing to write a lot of components ourselves and would make the site unless we worked really hard at it like it was written in a purely functional way without regard for aesthetics.

This project was last updated 2 years ago. At this point I'd consider the project either dead or unsupported.

##### Vugu

[Vugu](https://www.vugu.org/) is a direct competitor to Vecty. It is also Go based and experimental.

In a lot of ways the writing of the code looks like a cross between `tsx` and `go`, basically some html tags and go looking stuff.

There is an existing component [library](https://github.com/vugu/vugu/tree/master/vgform) for use in the project. However it is very very simple and covers so little I'm not sure it's of much added value. So much like Vecty we would be writing a lot of these ourselves. Vugu does have some [thoughts](https://github.com/vugu/vugu/wiki/Component-Library-Notes) on writing a library but is still very basic

Was last updated 3 weeks ago(June 2nd 2024), while not as recent as I'd like it is pretty good.

##### Go-App

[Go-App](https://github.com/maxence-charriere/go-app/) has recent development, but no existing component library.

It's declarative and intended to be composable, it's written in pure Go.

It looks like it's WASM and server-side rendering, which isn't strictly problematic but most people in the front-end space have been moving away from server-rendering mostly.

##### Vue

[Vue](https://github.com/norunners/vue) not to be confused by vue.js framework. It's intended to mimic the vue interface and style of writing which looks a bit like some weird html and go code in seperate files. It's an alpha project and not seen any updates since Dec 25, 2021.

---

#### Typescript/Javascript

None of the developers on this project care for TS/JS so anything written in that is starting at a disadvantage.

##### React

I have recent and somewhat extensive knowledge of working with [React](https://react.dev/), this would mean a simple but mostly okay to decent looking frontend is achievable in relatively short order.

using [Material UI](https://mui.com/material-ui/all-components/), [React i18n](https://react.i18next.com/), [React Router](https://reactrouter.com/en/main), [TanStack Query](https://tanstack.com/query/latest) would get us pretty far forward for fast launch.


##### Vue

I've not personally used [Vue](https://vuejs.org/) but I've heard good things(for TS), and would like to use it at some point but overall we would be starting mostly from scratch in terms of knowledge of getting it quickly stood up.

##### Angular

[Angular](https://angular.dev/)

Maintained by a dedicated team at Google, Angular provides a broad suite of tools, APIs, and libraries to simplify and streamline your development workflow. Angular gives you a solid platform on which to build fast, reliable applications that scale with both the size of your team and the size of your codebase.

##### NestJS

[NestJS](https://nestjs.com/) aims to be both frontend and backend.

Nest (NestJS) is a framework for building efficient, scalable Node.js server-side applications. It uses progressive JavaScript, is built with and fully supports TypeScript (yet still enables developers to code in pure JavaScript) and combines elements of OOP (Object Oriented Programming), FP (Functional Programming), and FRP (Functional Reactive Programming).

Nest provides an out-of-the-box application architecture which allows developers and teams to create highly testable, scalable, loosely coupled, and easily maintainable applications. The architecture is heavily inspired by Angular.

Under the hood, Nest makes use of robust HTTP Server frameworks like Express (the default) and optionally can be configured to use Fastify as well!

##### Svelte

[Svelte](https://svelte.dev/docs/introduction) is interesting, it implements scripts in js, that can be used in the `.svelte` files which are compiled to minimize the bundle size and increase performance. It's an interesting composition. I've not used it but have heard of it before.

---


## Decision

I'm not seeing anything stable enough to give me great confidence in it's ability to grant us the ability to deliver quickly a web ui with wasm. In light of this I would support using React, my experience with React and making the project less about the frontend and more about the backend, we should aim for simplicity on the frontend while keeping an eye on aesthetics.

## Consequences

The consequence of choosing React means additional package management tooling and code bloat. We will need a way to maintain types from the backend to the frontend, possibly using swagger tools to export our api contract to swagger for the frontend to build a client from.
