<template>
  <div>
    <template v-for="(col ,index) in cols">
      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'string')" :label="col.name">
        <el-input v-model="temp[col.name]" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'float')" :label="col.name">
        <el-input-number v-model="temp[col.name]" :precision="2" :step="0.01" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'double')" :label="col.name">
        <el-input-number v-model="temp[col.name]" :precision="2" :step="0.01" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'uint32')" :label="col.name">
        <el-input-number v-model="temp[col.name]" ontrols-position="right" :min="0" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'uint64')" :label="col.name">
        <el-input-number v-model="temp[col.name]" ontrols-position="right" :min="0" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'bool')" :label="col.name">
        <el-switch v-model="temp[col.name]" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'int64')" :label="col.name">
        <el-input-number v-model="temp[col.name]" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'int32')" :label="col.name">
        <el-input-number v-model="temp[col.name]" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'text')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'string_array')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'bool_array')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'uint32_array')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'uint64_array')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'float_array')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'double_array')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'int32_array')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'int64_array')" :label="col.name">
        <el-input v-model="temp[col.name]" type="textarea" :rows="3" :placeholder="col.desc" :disabled="checkDisable(col.feature)" />
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'enum')" :label="col.name">
        <el-select v-model="temp[col.name]">
          <el-option v-for="(field ,index) in col.data" :key="field" :label="field" :value="field" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'format')" :label="col.name">
        <el-select v-model="temp[col.name]" :disabled="checkDisable(col.feature)">
          <el-option v-for="(field ,index) in format_types" :key="field" :label="field" :value="field" />
        </el-select>
      </el-form-item>

      <el-form-item v-if="checkFeature(col.feature, 'edit') && (col.type === 'fields')" :label="col.name">
        <el-select v-model="temp[col.name]">
          <el-option v-for="(field ,index) in field_types" :key="field" :label="field" :value="field" />
        </el-select>
      </el-form-item>
    </template>
  </div>
</template>

<script>
  import store from '@/store'
  export default {
    name: 'TableForm',
    data() {
      return {
        features: store.getters.features,
        format_types: store.getters.format_types,
        field_types: store.getters.field_types,
      }
    },
    props: {
      cols: {},
      temp: {},
      dialogStatus: "",
    },
    methods: {
      checkDisable(feature) {
        if((feature & (1 << this.features['disableChange'])) > 0) {
          if(this.dialogStatus === "update") {
            return "disabled"
          }
        }
        return false
      },
      checkFeature(feature, name) {
        return (feature & (1 << this.features[name])) > 0
      },
    }
  }
</script>
