<template>
  <div class="dashboard-container">
    <panel-group />

  </div>
</template>

<script>
  import PanelGroup from './components/PanelGroup'

  export default {
    name: 'Dashboard',
    components: {
      PanelGroup,
    },
    data() {
      return {
      }
    },
    methods: {
    }
  }
</script>
