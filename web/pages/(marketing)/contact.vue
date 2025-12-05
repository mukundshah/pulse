<script setup lang="ts">
useHead({
  title: 'Contact Us',
})
</script>

<template>
  <div class="container mx-auto max-w-4xl py-16 px-4">
    <div class="mx-auto max-w-2xl text-center">
      <h1 class="mb-4 text-4xl font-semibold tracking-tight md:text-5xl">
        Get in touch
      </h1>
      <p class="mb-12 text-lg text-muted-foreground">
        Have a question or need help? We'd love to hear from you.
      </p>

      <div class="grid gap-8 md:grid-cols-2">
        <Card>
          <CardHeader>
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
              <Icon class="h-6 w-6" name="lucide:mail" />
            </div>
            <CardTitle>Email</CardTitle>
            <CardDescription>
              Send us an email and we'll get back to you as soon as possible.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <a
              class="text-foreground hover:text-foreground/80"
              href="mailto:support@pulse.app"
            >
              support@pulse.app
            </a>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-lg border border-border bg-background">
              <Icon class="h-6 w-6" name="lucide:message-circle" />
            </div>
            <CardTitle>Support</CardTitle>
            <CardDescription>
              Need technical assistance? Our support team is here to help.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <a
              class="text-foreground hover:text-foreground/80"
              href="mailto:help@pulse.app"
            >
              help@pulse.app
            </a>
          </CardContent>
        </Card>
      </div>

      <div class="mt-12 w-full text-left">
        <h2 class="mb-3 text-3xl font-semibold tracking-tight">
          Frequently Asked Questions
        </h2>
        <p class="mb-8 text-muted-foreground">
          Answers to common questions about Pulse and its features. If you have any other questions, please don't hesitate to contact us.
        </p>
        <Accordion
          collapsible
          class="w-full"
          default-value="response-time"
          type="single"
        >
          <AccordionItem value="response-time">
            <AccordionTrigger>
              How quickly do you respond?
            </AccordionTrigger>
            <AccordionContent>
              We typically respond to all inquiries within 24 hours during business days.
            </AccordionContent>
          </AccordionItem>
          <AccordionItem value="enterprise-support">
            <AccordionTrigger>
              Do you offer enterprise support?
            </AccordionTrigger>
            <AccordionContent>
              Yes, we offer dedicated support for enterprise customers. Contact us to learn more about our enterprise plans.
            </AccordionContent>
          </AccordionItem>
          <AccordionItem value="demo">
            <AccordionTrigger>
              Can I schedule a demo?
            </AccordionTrigger>
            <AccordionContent>
              Absolutely! Send us an email and we'll set up a time that works for you.
            </AccordionContent>
          </AccordionItem>
          <AccordionItem value="setup">
            <AccordionTrigger>
              How long does it take to set up?
            </AccordionTrigger>
            <AccordionContent>
              Setting up Pulse is quick and straightforward. Most users can get started in under 10 minutes by connecting their services and configuring basic monitoring rules.
            </AccordionContent>
          </AccordionItem>
          <AccordionItem value="data-security">
            <AccordionTrigger>
              Is my monitoring data secure?
            </AccordionTrigger>
            <AccordionContent>
              Yes, security is our top priority. All data is encrypted in transit and at rest. We follow industry best practices and comply with security standards to keep your information safe.
            </AccordionContent>
          </AccordionItem>
          <AccordionItem value="integrations">
            <AccordionTrigger>
              What services can I monitor?
            </AccordionTrigger>
            <AccordionContent>
              Pulse supports monitoring for a wide range of services including cloud providers (AWS, GCP, Azure), containers, databases, APIs, and custom applications. Check our integrations page for the full list.
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </div>
    </div>
  </div>
</template>
