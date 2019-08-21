
<template>
  <!-- eslint-disable vue/require-component-is -->
  <component v-bind="linkProps(to,params)">
    <slot />
  </component>
</template>

<script>
import { isExternal } from '@/utils/validate'

export default {
  props: {
    to: {
      type: String,
      required: true
    },
    params: {
      type: Object,
      required: false
    }
  },
  methods: {
    linkProps(url, params) {
      if (isExternal(url)) {
        return {
          is: 'a',
          href: url,
          target: '_blank',
          rel: 'noopener'
        }
      }
      return {
        is: 'router-link',
        to: {
          path: url,
          query: params
        }
      }
    }
  }
}
</script>
