<script setup lang="ts">
useHead({
  title: 'Terms of Service',
})
</script>

<template>
  <div class="container mx-auto max-w-4xl py-16 px-4">
    <div class="prose prose-invert max-w-none">
      <h1 class="mb-8 text-4xl font-semibold tracking-tight">
        Terms of Service
      </h1>
      <p class="text-muted-foreground">
        Last updated: <NuxtTime :datetime="new Date('2025-11-29T00:00:00Z').toISOString()" locale="en-US" v-bind="{ year: 'numeric', month: 'long', day: 'numeric' }"></NuxtTime>
      </p>

      <section class="mt-12">
        <h2>1. Acceptance of Terms</h2>
        <p>
          By accessing and using Pulse, you accept and agree to be bound by the terms and provision of this agreement.
        </p>
      </section>

      <section class="mt-8">
        <h2>2. Use License</h2>
        <p>
          Permission is granted to temporarily use Pulse for personal or commercial monitoring purposes. This is the grant of a license, not a transfer of title, and under this license you may not:
        </p>
        <ul>
          <li>Modify or copy the materials</li>
          <li>Use the materials for any commercial purpose without explicit written consent</li>
          <li>Attempt to reverse engineer any software contained in Pulse</li>
          <li>Remove any copyright or other proprietary notations from the materials</li>
        </ul>
      </section>

      <section class="mt-8">
        <h2>3. Service Availability</h2>
        <p>
          We strive to maintain high availability of our service, but we do not guarantee uninterrupted access. Pulse may be unavailable due to maintenance, updates, or circumstances beyond our control.
        </p>
      </section>

      <section class="mt-8">
        <h2>4. User Responsibilities</h2>
        <p>
          You are responsible for:
        </p>
        <ul>
          <li>Maintaining the confidentiality of your account credentials</li>
          <li>All activities that occur under your account</li>
          <li>Ensuring compliance with all applicable laws and regulations</li>
          <li>Not using the service for any illegal or unauthorized purpose</li>
        </ul>
      </section>

      <section class="mt-8">
        <h2>5. Data and Privacy</h2>
        <p>
          Your use of Pulse is also governed by our Privacy Policy. Please review our Privacy Policy to understand our practices regarding data collection and use.
        </p>
      </section>

      <section class="mt-8">
        <h2>6. Limitation of Liability</h2>
        <p>
          In no event shall Pulse or its suppliers be liable for any damages (including, without limitation, damages for loss of data or profit, or due to business interruption) arising out of the use or inability to use Pulse, even if Pulse or a Pulse authorized representative has been notified orally or in writing of the possibility of such damage.
        </p>
      </section>

      <section class="mt-8">
        <h2>7. Termination</h2>
        <p>
          We may terminate or suspend your account and access to the service immediately, without prior notice or liability, for any reason whatsoever, including without limitation if you breach the Terms.
        </p>
      </section>

      <section class="mt-8">
        <h2>8. Changes to Terms</h2>
        <p>
          We reserve the right, at our sole discretion, to modify or replace these Terms at any time. If a revision is material, we will provide at least 30 days notice prior to any new terms taking effect.
        </p>
      </section>

      <section class="mt-8">
        <h2>9. Contact Information</h2>
        <p>
          If you have any questions about these Terms, please contact us at <a href="/contact">contact</a>.
        </p>
      </section>
    </div>
  </div>
</template>

<style scoped>
@reference 'assets/css/tailwind.css';

.prose {
  @apply text-foreground;
}

.prose h1 {
  @apply mb-8 text-4xl font-semibold tracking-tight text-foreground;
}

.prose h2 {
  @apply mt-8 mb-4 text-2xl font-semibold text-foreground;
}

.prose h3 {
  @apply mt-6 mb-3 text-xl font-semibold text-foreground;
}

.prose p {
  @apply mb-4 text-muted-foreground leading-relaxed;
}

.prose ul,
.prose ol {
  @apply mb-4 ml-6 text-muted-foreground;
}

.prose li {
  @apply mb-2;
}

.prose a {
  @apply text-foreground underline underline-offset-4 hover:text-foreground/80;
}
</style>
