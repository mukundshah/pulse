<script setup lang="ts">
useHead({
  title: 'Privacy Policy',
})
</script>

<template>
  <div class="container mx-auto max-w-4xl py-16 px-4">
    <div class="prose prose-invert max-w-none">
      <h1 class="mb-8 text-4xl font-semibold tracking-tight">
        Privacy Policy
      </h1>
      <p class="text-muted-foreground">
        Last updated: <NuxtTime :datetime="new Date('2025-11-29T00:00:00Z').toISOString()" locale="en-US" v-bind="{ year: 'numeric', month: 'long', day: 'numeric' }"></NuxtTime>
      </p>

      <section class="mt-12">
        <h2>1. Introduction</h2>
        <p>
          Pulse ("we", "our", or "us") is committed to protecting your privacy. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you use our monitoring service.
        </p>
      </section>

      <section class="mt-8">
        <h2>2. Information We Collect</h2>
        <h3>2.1 Information You Provide</h3>
        <p>
          We collect information that you provide directly to us, including:
        </p>
        <ul>
          <li>Account registration information (name, email address)</li>
          <li>Payment information (processed through secure third-party providers)</li>
          <li>Communication data when you contact us</li>
        </ul>

        <h3>2.2 Automatically Collected Information</h3>
        <p>
          When you use Pulse, we automatically collect certain information, including:
        </p>
        <ul>
          <li>Usage data and analytics</li>
          <li>Device information and IP addresses</li>
          <li>Cookies and similar tracking technologies</li>
        </ul>

        <h3>2.3 Monitoring Data</h3>
        <p>
          We collect and process monitoring data from the services and infrastructure you connect to Pulse. This data is used solely for providing monitoring and alerting services.
        </p>
      </section>

      <section class="mt-8">
        <h2>3. How We Use Your Information</h2>
        <p>
          We use the information we collect to:
        </p>
        <ul>
          <li>Provide, maintain, and improve our services</li>
          <li>Process transactions and send related information</li>
          <li>Send technical notices, updates, and support messages</li>
          <li>Respond to your comments and questions</li>
          <li>Monitor and analyze trends and usage</li>
          <li>Detect, prevent, and address technical issues</li>
        </ul>
      </section>

      <section class="mt-8">
        <h2>4. Data Storage and Security</h2>
        <p>
          We implement appropriate technical and organizational measures to protect your personal information. However, no method of transmission over the Internet or electronic storage is 100% secure.
        </p>
        <p>
          Your monitoring data is encrypted in transit and at rest. We retain your data only for as long as necessary to provide our services and comply with legal obligations.
        </p>
      </section>

      <section class="mt-8">
        <h2>5. Data Sharing and Disclosure</h2>
        <p>
          We do not sell, trade, or rent your personal information to third parties. We may share your information only in the following circumstances:
        </p>
        <ul>
          <li>With service providers who assist us in operating our service</li>
          <li>When required by law or to protect our rights</li>
          <li>In connection with a business transfer or merger</li>
          <li>With your explicit consent</li>
        </ul>
      </section>

      <section class="mt-8">
        <h2>6. Your Rights</h2>
        <p>
          Depending on your location, you may have certain rights regarding your personal information, including:
        </p>
        <ul>
          <li>The right to access your personal information</li>
          <li>The right to correct inaccurate information</li>
          <li>The right to delete your personal information</li>
          <li>The right to object to processing of your information</li>
          <li>The right to data portability</li>
        </ul>
        <p>
          To exercise these rights, please contact us through our <a href="/contact">contact page</a>.
        </p>
      </section>

      <section class="mt-8">
        <h2>7. Cookies and Tracking Technologies</h2>
        <p>
          We use cookies and similar tracking technologies to track activity on our service and hold certain information. You can instruct your browser to refuse all cookies or to indicate when a cookie is being sent.
        </p>
      </section>

      <section class="mt-8">
        <h2>8. Children's Privacy</h2>
        <p>
          Our service is not intended for individuals under the age of 18. We do not knowingly collect personal information from children.
        </p>
      </section>

      <section class="mt-8">
        <h2>9. Changes to This Privacy Policy</h2>
        <p>
          We may update our Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on this page and updating the "Last updated" date.
        </p>
      </section>

      <section class="mt-8">
        <h2>10. Contact Us</h2>
        <p>
          If you have any questions about this Privacy Policy, please contact us at <a href="/contact">contact</a>.
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
