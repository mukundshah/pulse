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
        <div class="mb-8 inline-flex items-center gap-2 rounded-full border border-border bg-muted/50 px-4 py-1.5 text-sm backdrop-blur-sm">
          <span class="relative flex h-2 w-2">
            <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-foreground opacity-75"></span>
            <span class="relative inline-flex h-2 w-2 rounded-full bg-foreground"></span>
          </span>
          <span class="text-muted-foreground">Now in public beta</span>
        </div>
        <h1 class="mb-6 text-5xl font-semibold tracking-tight md:text-6xl lg:text-7xl">
          The modern way to
          <span class="relative inline-block">
            <span class="bg-linear-to-r from-foreground via-foreground/80 to-muted-foreground bg-clip-text text-transparent">
              monitor
            </span>
            <span class="absolute -bottom-2 left-0 h-1 w-full bg-linear-to-r from-foreground/20 via-foreground/10 to-transparent"></span>
          </span>
          your infrastructure
        </h1>
        <p class="mb-10 text-lg text-muted-foreground md:text-xl">
          Pulse provides real-time monitoring, alerting, and insights for your applications and services.
          Built for teams who care about reliability.
        </p>
        <div class="flex flex-col items-center justify-center gap-4 sm:flex-row">
          <Button as-child size="lg">
            <NuxtLink to="/auth/login">
              Get started
            </NuxtLink>
          </Button>
          <Button
            as-child
            size="lg"
            variant="outline"
          >
            <a href="#features">
              Learn more
            </a>
          </Button>
        </div>

        <!-- Stats preview -->
        <div class="mt-16 grid grid-cols-3 gap-8 border-t border-border/50 pt-12">
          <div class="flex flex-col items-center">
            <div class="mb-2 text-3xl font-semibold md:text-4xl">
              99.9%
            </div>
            <div class="text-sm text-muted-foreground">
              Uptime SLA
            </div>
          </div>
          <div class="flex flex-col items-center">
            <div class="mb-2 text-3xl font-semibold md:text-4xl">
              &lt;100ms
            </div>
            <div class="text-sm text-muted-foreground">
              Response Time
            </div>
          </div>
          <div class="flex flex-col items-center">
            <div class="mb-2 text-3xl font-semibold md:text-4xl">
              24/7
            </div>
            <div class="text-sm text-muted-foreground">
              Monitoring
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="relative border-t border-border bg-muted/30 py-24">
      <div class="container relative z-10">
        <div class="mb-16 text-center">
          <div class="mb-4 inline-flex items-center gap-2">
            <div class="h-px w-8 bg-border"></div>
            <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Features</span>
            <div class="h-px w-8 bg-border"></div>
          </div>
          <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
            Everything you need to stay on top
          </h2>
          <p class="mx-auto max-w-2xl text-muted-foreground">
            Powerful features designed to help you monitor, alert, and respond to issues faster.
          </p>
        </div>
        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <div class="absolute -right-4 -top-4 h-24 w-24 rounded-full bg-foreground/5 blur-2xl transition-all duration-500 group-hover:scale-150"></div>
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background transition-all duration-300 group-hover:scale-110 group-hover:border-foreground/30">
                <Icon class="h-6 w-6 transition-transform duration-300 group-hover:scale-110" name="lucide:activity" />
              </div>
              <CardTitle>Real-time Monitoring</CardTitle>
              <CardDescription>
                Monitor your services in real-time with instant updates and comprehensive metrics.
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <div class="absolute -right-4 -top-4 h-24 w-24 rounded-full bg-foreground/5 blur-2xl transition-all duration-500 group-hover:scale-150"></div>
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background transition-all duration-300 group-hover:scale-110 group-hover:border-foreground/30">
                <Icon class="h-6 w-6 transition-transform duration-300 group-hover:scale-110" name="lucide:bell" />
              </div>
              <CardTitle>Smart Alerting</CardTitle>
              <CardDescription>
                Get notified instantly when something goes wrong with intelligent alert routing.
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <div class="absolute -right-4 -top-4 h-24 w-24 rounded-full bg-foreground/5 blur-2xl transition-all duration-500 group-hover:scale-150"></div>
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background transition-all duration-300 group-hover:scale-110 group-hover:border-foreground/30">
                <Icon class="h-6 w-6 transition-transform duration-300 group-hover:scale-110" name="lucide:trending-up" />
              </div>
              <CardTitle>Performance Insights</CardTitle>
              <CardDescription>
                Understand your system's performance with detailed analytics and trends.
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <div class="absolute -right-4 -top-4 h-24 w-24 rounded-full bg-foreground/5 blur-2xl transition-all duration-500 group-hover:scale-150"></div>
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background transition-all duration-300 group-hover:scale-110 group-hover:border-foreground/30">
                <Icon class="h-6 w-6 transition-transform duration-300 group-hover:scale-110" name="lucide:shield" />
              </div>
              <CardTitle>Secure & Reliable</CardTitle>
              <CardDescription>
                Enterprise-grade security with 99.9% uptime SLA and data encryption.
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <div class="absolute -right-4 -top-4 h-24 w-24 rounded-full bg-foreground/5 blur-2xl transition-all duration-500 group-hover:scale-150"></div>
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background transition-all duration-300 group-hover:scale-110 group-hover:border-foreground/30">
                <Icon class="h-6 w-6 transition-transform duration-300 group-hover:scale-110" name="lucide:zap" />
              </div>
              <CardTitle>Lightning Fast</CardTitle>
              <CardDescription>
                Built for speed with sub-second response times and efficient resource usage.
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <div class="absolute -right-4 -top-4 h-24 w-24 rounded-full bg-foreground/5 blur-2xl transition-all duration-500 group-hover:scale-150"></div>
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background transition-all duration-300 group-hover:scale-110 group-hover:border-foreground/30">
                <Icon class="h-6 w-6 transition-transform duration-300 group-hover:scale-110" name="lucide:users" />
              </div>
              <CardTitle>Team Collaboration</CardTitle>
              <CardDescription>
                Work together seamlessly with shared dashboards and team notifications.
              </CardDescription>
            </CardHeader>
          </Card>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section id="how-it-works" class="relative border-t border-border py-24">
      <div class="container relative z-10">
        <div class="mb-16 text-center">
          <div class="mb-4 inline-flex items-center gap-2">
            <div class="h-px w-8 bg-border"></div>
            <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">How it works</span>
            <div class="h-px w-8 bg-border"></div>
          </div>
          <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
            Simple setup, powerful results
          </h2>
          <p class="mx-auto max-w-2xl text-muted-foreground">
            Get started in minutes with our straightforward setup process.
          </p>
        </div>
        <div class="grid gap-8 md:grid-cols-3">
          <div class="relative flex flex-col items-center text-center">
            <div class="mb-6 flex h-16 w-16 items-center justify-center rounded-full border border-border bg-muted/50">
              <Icon class="h-8 w-8" name="lucide:plug" />
            </div>
            <h3 class="mb-2 text-xl font-semibold">
              Connect
            </h3>
            <p class="text-muted-foreground">
              Connect your services and infrastructure with our simple integration process. No complex configuration needed.
            </p>
          </div>
          <div class="relative flex flex-col items-center text-center">
            <div class="mb-6 flex h-16 w-16 items-center justify-center rounded-full border border-border bg-muted/50">
              <Icon class="h-8 w-8" name="lucide:settings" />
            </div>
            <h3 class="mb-2 text-xl font-semibold">
              Configure
            </h3>
            <p class="text-muted-foreground">
              Set up monitoring rules, alert thresholds, and notification preferences in our intuitive dashboard.
            </p>
          </div>
          <div class="relative flex flex-col items-center text-center">
            <div class="mb-6 flex h-16 w-16 items-center justify-center rounded-full border border-border bg-muted/50">
              <Icon class="h-8 w-8" name="lucide:eye" />
            </div>
            <h3 class="mb-2 text-xl font-semibold">
              Monitor
            </h3>
            <p class="text-muted-foreground">
              Watch your infrastructure in real-time with comprehensive dashboards and instant alerts.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Integrations Section -->
    <section id="integrations" class="relative border-t border-border bg-muted/30 py-24">
      <div class="container relative z-10">
        <div class="mb-16 text-center">
          <div class="mb-4 inline-flex items-center gap-2">
            <div class="h-px w-8 bg-border"></div>
            <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Integrations</span>
            <div class="h-px w-8 bg-border"></div>
          </div>
          <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
            Works with your stack
          </h2>
          <p class="mx-auto max-w-2xl text-muted-foreground">
            Integrate with the tools and services you already use.
          </p>
        </div>
        <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20">
            <CardHeader class="relative flex flex-col items-center text-center">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:cloud" />
              </div>
              <CardTitle class="text-base">
                Cloud Providers
              </CardTitle>
              <CardDescription class="text-sm">
                AWS, GCP, Azure
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20">
            <CardHeader class="relative flex flex-col items-center text-center">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:server" />
              </div>
              <CardTitle class="text-base">
                Containers
              </CardTitle>
              <CardDescription class="text-sm">
                Docker, Kubernetes
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20">
            <CardHeader class="relative flex flex-col items-center text-center">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:code" />
              </div>
              <CardTitle class="text-base">
                APIs
              </CardTitle>
              <CardDescription class="text-sm">
                REST, GraphQL
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20">
            <CardHeader class="relative flex flex-col items-center text-center">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:database" />
              </div>
              <CardTitle class="text-base">
                Databases
              </CardTitle>
              <CardDescription class="text-sm">
                PostgreSQL, MySQL, MongoDB
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20">
            <CardHeader class="relative flex flex-col items-center text-center">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:message-square" />
              </div>
              <CardTitle class="text-base">
                Notifications
              </CardTitle>
              <CardDescription class="text-sm">
                Slack, Email, PagerDuty
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20">
            <CardHeader class="relative flex flex-col items-center text-center">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:git-branch" />
              </div>
              <CardTitle class="text-base">
                CI/CD
              </CardTitle>
              <CardDescription class="text-sm">
                GitHub Actions, GitLab
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20">
            <CardHeader class="relative flex flex-col items-center text-center">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:webhook" />
              </div>
              <CardTitle class="text-base">
                Webhooks
              </CardTitle>
              <CardDescription class="text-sm">
                Custom integrations
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20">
            <CardHeader class="relative flex flex-col items-center text-center">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:plus" />
              </div>
              <CardTitle class="text-base">
                More
              </CardTitle>
              <CardDescription class="text-sm">
                And many more...
              </CardDescription>
            </CardHeader>
          </Card>
        </div>
      </div>
    </section>

    <!-- Use Cases Section -->
    <section class="relative border-t border-border py-24">
      <div class="container relative z-10">
        <div class="mb-16 text-center">
          <div class="mb-4 inline-flex items-center gap-2">
            <div class="h-px w-8 bg-border"></div>
            <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Use cases</span>
            <div class="h-px w-8 bg-border"></div>
          </div>
          <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
            Built for every team
          </h2>
          <p class="mx-auto max-w-2xl text-muted-foreground">
            Whether you're a startup or enterprise, Pulse scales with you.
          </p>
        </div>
        <div class="grid gap-6 md:grid-cols-2">
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:rocket" />
              </div>
              <CardTitle>Startups & Scale-ups</CardTitle>
              <CardDescription>
                Get started quickly with minimal setup. Monitor your growing infrastructure without the complexity of enterprise tools.
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:building-2" />
              </div>
              <CardTitle>Enterprise Teams</CardTitle>
              <CardDescription>
                Scale to thousands of services with advanced features, role-based access, and comprehensive audit logs.
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:code-2" />
              </div>
              <CardTitle>Development Teams</CardTitle>
              <CardDescription>
                Integrate monitoring into your workflow with CI/CD pipelines and development-friendly APIs.
              </CardDescription>
            </CardHeader>
          </Card>
          <Card class="group relative overflow-hidden transition-all duration-300 hover:border-foreground/20 hover:shadow-lg">
            <CardHeader class="relative">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
                <Icon class="h-6 w-6" name="lucide:shield-check" />
              </div>
              <CardTitle>DevOps & SRE</CardTitle>
              <CardDescription>
                Advanced alerting, incident management, and reliability metrics to keep your systems running smoothly.
              </CardDescription>
            </CardHeader>
          </Card>
        </div>
      </div>
    </section>

    <!-- Visual separator -->
    <div class="relative border-t border-border">
      <div class="absolute left-1/2 top-0 h-px w-32 -translate-x-1/2 -translate-y-1/2 bg-gradient-to-r from-transparent via-border to-transparent"></div>
    </div>

    <!-- CTA Section -->
    <section class="group relative border-t border-border py-24">
      <!-- Background decoration - blended gradients -->
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-foreground/5 via-transparent to-transparent"></div>
      <div class="absolute inset-0 bg-gradient-to-br from-foreground/5 via-transparent to-foreground/5 transition-opacity duration-500 group-hover:opacity-100"></div>

      <div class="relative z-10 container flex flex-col items-center justify-center text-center">
        <div class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-full border border-border bg-muted/50">
          <Icon class="h-6 w-6" name="lucide:rocket" />
        </div>
        <h2 class="mb-4 text-3xl font-semibold tracking-tight md:text-4xl">
          Ready to get started?
        </h2>
        <p class="mb-10 mx-auto max-w-2xl text-base text-muted-foreground md:text-lg">
          Join teams already using Pulse to monitor their infrastructure.
        </p>
        <div class="flex flex-col items-center justify-center gap-4 sm:flex-row">
          <Button as-child size="lg">
            <NuxtLink to="/auth/login">
              Start monitoring
            </NuxtLink>
          </Button>
          <Button as-child size="lg" variant="outline">
            <a href="#features">View features</a>
          </Button>
        </div>
      </div>
    </section>
  </div>
</template>
