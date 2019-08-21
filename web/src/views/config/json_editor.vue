<template>
  <div class="createPost-container">
    <sticky :z-index="10" :class-name="'sub-navbar draft'">
      <el-button style="margin-left: 10px;" type="success" @click="submit">
        Publush
      </el-button>
    </sticky>

    <div class="components-container">
      <div class="editor-container">
        <json-editor ref="jsonEditor" v-model="value" />
      </div>
    </div>
  </div>
</template>

<script>
import Sticky from '@/components/Sticky' // 粘性header组件
import JsonEditor from '@/components/JsonEditor'
import { getData, setData } from '@/api/config_json'

export default {
  name: 'JsonEditorDemo',
  components: { JsonEditor, Sticky },
  data() {
    return {
      value: JSON.parse('{}'),
      query: {}
    }
  },
  mounted() {
    this.render()
  },
  methods: {
    render() {
      this.query = this.$route.query
      getData(this.query).then(response => {
        this.value = JSON.parse(response.data)
      })
    },
    submit() {
      setData(this.query, this.value).then(response => {
        this.$notify({
          title: 'Success',
          message: 'Publish Success',
          type: 'success',
          duration: 2000
        })
      })
    }
  }
}
</script>

<style scoped>
.editor-container{
  position: relative;
  height: 100%;
}
</style>

