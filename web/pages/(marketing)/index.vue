<template>
  <div class="relative">
    <!-- Background decorative elements -->
    <div class="pointer-events-none fixed inset-0 overflow-hidden">
      <!-- Gradient orbs -->
      <div class="absolute -left-1/4 top-0 h-[600px] w-[600px] rounded-full bg-foreground/5 blur-3xl"></div>
      <div class="absolute -right-1/4 top-1/3 h-[500px] w-[500px] rounded-full bg-foreground/5 blur-3xl"></div>
      <div class="absolute bottom-0 left-1/2 h-[400px] w-[400px] -translate-x-1/2 rounded-full bg-foreground/5 blur-3xl"></div>
    </div>

    <!-- Hero Section -->
    <section class="relative container flex flex-1 flex-col items-center justify-center py-24 md:py-32">
      <div class="relative z-10 mx-auto max-w-4xl text-center">
        <h1 class="mb-6 text-5xl font-semibold tracking-tight md:text-6xl lg:text-7xl">
          Know before
          <span class="relative inline-block">
            <span class="bg-linear-to-r from-foreground via-foreground/80 to-muted-foreground bg-clip-text text-transparent">
              your users
            </span>
            <span class="absolute -bottom-2 left-0 h-1 w-full bg-linear-to-r from-foreground/20 via-foreground/10 to-transparent"></span>
          </span>
          do
        </h1>
        <p class="mb-10 text-lg text-muted-foreground md:text-xl">
          Pulse is a self-hosted synthetic monitoring platform built from the ground up.
          Born from curiosity about how monitoring systems work, crafted with care to understand every detail.
        </p>
        <div class="flex flex-col items-center justify-center gap-4 sm:flex-row">
          <Button as-child size="lg">
            <NuxtLink to="/auth/register">
              Get started
            </NuxtLink>
          </Button>
          <Button
            as-child
            size="lg"
            variant="outline"
          >
            <a href="#story">
              Learn the story
            </a>
          </Button>
        </div>
      </div>
    </section>

    <!-- The Story Section -->
    <section id="story" class="relative border-t border-border py-24">
      <div class="container relative z-10">
        <div class="mx-auto max-w-3xl">
          <div class="mb-16 text-center">
            <div class="mb-4 inline-flex items-center gap-2">
              <div class="h-px w-8 bg-border"></div>
              <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">The Story</span>
              <div class="h-px w-8 bg-border"></div>
            </div>
            <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
              Built from curiosity, not necessity
            </h2>
          </div>

          <div class="space-y-8 text-lg leading-relaxed text-muted-foreground">
            <p>
              Every great product starts with a question. For Pulse, that question was simple: <em class="text-foreground">How do monitoring platforms actually work?</em>
            </p>
            <p>
              Inspired by tools like Checkly, I wanted to understand not just what synthetic monitoring does, but how it's built. What are the architectural decisions? How do you handle millions of checks? What makes a monitoring system reliable enough to trust with your infrastructure?
            </p>
            <p>
              So I built Pulse. Not because I needed another monitoring tool, but because I needed to understand. Every line of code, every design decision, every feature was a lesson in building systems that matter.
            </p>
            <p class="text-foreground font-medium">
              This is what happens when curiosity meets dedication. When you're not satisfied with just using a tool, but need to understand how it works from the inside out.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- What It Does Section -->
    <section id="features" class="relative border-t border-border bg-muted/30 py-24">
      <div class="container relative z-10">
        <div class="mb-16 text-center">
          <div class="mb-4 inline-flex items-center gap-2">
            <div class="h-px w-8 bg-border"></div>
            <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">What It Does</span>
            <div class="h-px w-8 bg-border"></div>
          </div>
          <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
            Monitor what matters, know what's broken
          </h2>
          <p class="mx-auto max-w-2xl text-muted-foreground">
            Pulse watches your services around the clock, from multiple locations, and tells you the moment something goes wrong.
          </p>
        </div>
        <div class="grid gap-8 md:grid-cols-3">
          <div class="flex flex-col">
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
              <Icon class="h-6 w-6" name="lucide:globe" />
            </div>
            <h3 class="mb-2 text-xl font-semibold">
              HTTP & API Monitoring
            </h3>
            <p class="text-muted-foreground">
              Watch your APIs and web services. Know instantly when they're down, slow, or returning errors. Verify responses, check headers, validate that everything works as expected.
            </p>
          </div>
          <div class="flex flex-col">
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
              <Icon class="h-6 w-6" name="lucide:network" />
            </div>
            <h3 class="mb-2 text-xl font-semibold">
              Infrastructure Checks
            </h3>
            <p class="text-muted-foreground">
              Monitor TCP connections, DNS resolution, and network services. Ensure your databases, caches, and critical infrastructure are reachable and responding.
            </p>
          </div>
          <div class="flex flex-col">
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
              <Icon class="h-6 w-6" name="lucide:bell" />
            </div>
            <h3 class="mb-2 text-xl font-semibold">
              Smart Alerting
            </h3>
            <p class="text-muted-foreground">
              Get notified the moment something changes. Email alerts, status tracking, and a complete history of every incident. Know what happened, when, and why.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- The Philosophy Section -->
    <section class="relative border-t border-border py-24">
      <div class="container relative z-10">
        <div class="mx-auto max-w-3xl">
          <div class="mb-16 text-center">
            <div class="mb-4 inline-flex items-center gap-2">
              <div class="h-px w-8 bg-border"></div>
              <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">The Philosophy</span>
              <div class="h-px w-8 bg-border"></div>
            </div>
            <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
              Self-hosted by design
            </h2>
          </div>

          <div class="space-y-6 text-lg leading-relaxed text-muted-foreground">
            <p>
              Pulse is built to run on your infrastructure. No vendor lock-in, no monthly fees, no limits on what you can monitor. You own your data, you control your monitoring.
            </p>
            <p>
              This isn't about competing with commercial solutions—it's about understanding. About having the freedom to see exactly how everything works, to modify it to your needs, to learn from it.
            </p>
            <p class="text-foreground font-medium">
              Sometimes the best way to understand a system is to build it yourself. Pulse is that journey, shared.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- The Dedication Section -->
    <section class="relative border-t border-border bg-muted/30 py-24">
      <div class="container relative z-10">
        <div class="mx-auto max-w-3xl">
          <div class="mb-16 text-center">
            <div class="mb-4 inline-flex items-center gap-2">
              <div class="h-px w-8 bg-border"></div>
              <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">The Dedication</span>
              <div class="h-px w-8 bg-border"></div>
            </div>
            <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
              Every detail matters
            </h2>
          </div>

          <div class="grid gap-8 md:grid-cols-2">
            <Card class="border-foreground/10">
              <CardHeader>
                <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                  <Icon class="h-6 w-6" name="lucide:target" />
                </div>
                <CardTitle>Thoughtful Design</CardTitle>
                <CardDescription class="mt-2">
                  Every feature exists because it solves a real problem. No bloat, no unnecessary complexity. Just what you need to monitor effectively.
                </CardDescription>
              </CardHeader>
            </Card>
            <Card class="border-foreground/10">
              <CardHeader>
                <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                  <Icon class="h-6 w-6" name="lucide:code-2" />
                </div>
                <CardTitle>Built to Learn</CardTitle>
                <CardDescription class="mt-2">
                  The codebase is clean, documented, and structured to be understood. This is as much a learning resource as it is a tool.
                </CardDescription>
              </CardHeader>
            </Card>
            <Card class="border-foreground/10">
              <CardHeader>
                <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                  <Icon class="h-6 w-6" name="lucide:users" />
                </div>
                <CardTitle>Team Ready</CardTitle>
                <CardDescription class="mt-2">
                  Projects, teams, permissions, and collaboration built in from the start. This isn't a solo tool—it's designed for real teams.
                </CardDescription>
              </CardHeader>
            </Card>
            <Card class="border-foreground/10">
              <CardHeader>
                <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                  <Icon class="h-6 w-6" name="lucide:zap" />
                </div>
                <CardTitle>Real-time Everything</CardTitle>
                <CardDescription class="mt-2">
                  See check results as they happen. No refreshing, no polling. When something changes, you know immediately.
                </CardDescription>
              </CardHeader>
            </Card>
          </div>
        </div>
      </div>
    </section>

    <!-- The Vision Section -->
    <section class="relative border-t border-border py-24">
      <div class="container relative z-10">
        <div class="mx-auto max-w-3xl text-center">
          <div class="mb-8">
            <div class="mb-4 inline-flex items-center gap-2">
              <div class="h-px w-8 bg-border"></div>
              <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">The Vision</span>
              <div class="h-px w-8 bg-border"></div>
            </div>
            <h2 class="mb-6 text-3xl font-semibold tracking-tight md:text-4xl">
              Not a product, a journey
            </h2>
          </div>

          <div class="space-y-6 text-lg leading-relaxed text-muted-foreground">
            <p>
              Pulse isn't trying to be the next unicorn startup. It's not trying to disrupt an industry or raise venture capital.
            </p>
            <p>
              It's a project born from genuine curiosity about how complex systems work. It's built with care, attention to detail, and a deep respect for the craft of software engineering.
            </p>
            <p class="text-foreground font-medium">
              If you're someone who values understanding over convenience, who wants to know how things work under the hood, who believes in self-hosting and open source—this is for you.
            </p>
            <p>
              Try it. Learn from it. Contribute to it. Or just appreciate the dedication that went into building something meaningful, one commit at a time.
            </p>
          </div>

          <div class="mt-12 flex flex-col items-center justify-center gap-4 sm:flex-row">
            <Button as-child size="lg">
              <NuxtLink to="/auth/register">
                Start monitoring
              </NuxtLink>
            </Button>
            <Button as-child size="lg" variant="outline">
              <a href="https://github.com/mukundshah/pulse" rel="noopener noreferrer" target="_blank">
                <Icon class="mr-2 h-4 w-4" name="lucide:github" />
                Explore the code
              </a>
            </Button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
