<template>
  <div class="app-container">
    <div class="filter-container">

      <template v-for="(col ,index) in cols">
        <el-input v-if="checkFeature(col.feature, 'search')" v-model="listQuery[col.name]" :placeholder="col.name" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      </template>

      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        {{ $t('table.search') }}
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        {{ $t('table.add') }}
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-success" @click="handlePush">
        {{ $t('table.push') }}
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-plus" @click="handleMerge">
        {{ $t('table.merge') }}
      </el-button>
      <br>
      <a target="_blank" :href="exportUrl+'/excel'">
        <el-button v-waves :loading="downloadLoading" class="filter-item" type="primary" icon="el-icon-download">
          {{ $t('table.exportExcel') }}
        </el-button>
      </a>
      <a target="_blank" :href="exportUrl+'/json'">
        <el-button v-waves :loading="downloadLoading" class="filter-item" type="primary" icon="el-icon-download">
          {{ $t('table.exportJson') }}
        </el-button>
      </a>
      <a target="_blank" :href="exportUrl+'/lua'">
        <el-button v-waves :loading="downloadLoading" class="filter-item" type="primary" icon="el-icon-download">
          {{ $t('table.exportLua') }}
        </el-button>
      </a>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
      @sort-change="sortChange"
    >
      <template v-for="(col ,index) in cols">
        <el-table-column :prop="col.name" :sortable="checkFeature(col.feature, 'sort') ? 'custom' : false" :label="col.name" align="center">
          <template slot-scope="scope">
            <span>{{ scope.row.values[col.name] }}</span>
          </template>
        </el-table-column>
      </template>

      <el-table-column :label="$t('table.actions')" align="center" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            {{ $t('table.edit') }}
          </el-button>
          <el-button v-if="row.status!='deleted'" size="mini" type="danger" @click="handleDelete(row)">
            {{ $t('table.delete') }}
          </el-button>
        </template>
      </el-table-column>

    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="render" />

    <upload-excel-component :on-success="handleSuccess" :before-upload="beforeUpload" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="150px" style="width: 400px; margin-left:50px;">

        <table-form :cols="cols" :dialog-status="dialogStatus" :temp="temp" />

      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          {{ $t('table.cancel') }}
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          {{ $t('table.confirm') }}
        </el-button>
      </div>
    </el-dialog>

    <el-dialog :title="textMap[dialogConfirmStatus]" :visible.sync="dialogConfirmVisible">
      <el-alert
        v-if="dialogConfirmStatus==='delete'"
        :title="textMap['confirmMsg']"
        type="warning"
        effect="dark"
      />

      <template v-if="dialogConfirmStatus==='push'">
        <el-radio v-for="item in servers" v-model="pushServer" :label="item.id" size="large" border>{{ item.name }}</el-radio>
      </template>

      <template v-if="dialogConfirmStatus==='merge'">
        <el-timeline>
          <el-timeline-item
            v-for="(activity, index) in activities"
            :key="index"
            :color="activity.color"
            :timestamp="activity.timestamp"
          >
            {{ activity.content }}
          </el-timeline-item>
        </el-timeline>

        <el-alert
          v-if="dialogConfirmButton===false"
          :title="merge_error"
          type="warning"
          effect="dark"
        />
      </template>

      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogConfirmVisible = false">
          {{ $t('table.cancel') }}
        </el-button>
        <el-button v-if="dialogConfirmButton" type="primary" @click="handleConfirm">
          {{ $t('table.confirm') }}
        </el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { fetchList, createRow, updateRow, deleteRow, uploadData } from '@/api/config_table'
import { pushConfig, mergeConfigInfo, mergeConfig } from '@/api/config'
import waves from '@/directive/waves' // waves directive
import store from '@/store'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import TableForm from '@/components/TableForm'
import UploadExcelComponent from '@/components/UploadFile/excel.vue'

export default {
  name: 'ComplexTable',
  components: { Pagination, TableForm, UploadExcelComponent },
  directives: { waves },
  filters: {},
  data() {
    return {
      servers: store.getters.userSelectServers,
      tableKey: 0,
      cols: null,
      list: null,
      total: 0,
      exportUrl: '',
      activities: [],
      merge_error: '',
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        env_id: '',
        collect: '',
        sort: 'id'
      },
      temp: {},
      tempId: '',
      dialogFormVisible: false,
      dialogStatus: '',
      features: store.getters.features,
      textMap: {
        update: 'Edit',
        create: 'Create',
        delete: 'Confirm',
        push: 'Select a server',
        merge: 'History',
        confirmMsg: 'Are you sure?'
      },
      pushServer: '',
      dialogConfirmVisible: false,
      dialogConfirmButton: true,
      dialogConfirmStatus: '',
      deleteRow: {},
      tHeader: [],
      filterVal: [],
      rules: {
        type: [{ required: true, message: 'type is required', trigger: 'change' }],
        timestamp: [{ type: 'date', required: true, message: 'timestamp is required', trigger: 'change' }],
        title: [{ required: true, message: 'title is required', trigger: 'blur' }]
      },
      downloadLoading: false
    }
  },
  mounted() {
    this.render()
  },
  methods: {
    render() {
      this.listLoading = true
      this.listQuery.env_id = this.$route.query.envId
      this.listQuery.collect = this.$route.name
      this.exportUrl = '/v1/table/' + this.listQuery.env_id + '/' + this.listQuery.collect
      fetchList(this.listQuery).then(response => {
        this.list = response.data.items
        this.cols = response.data.cols
        this.total = response.data.total
        for (let i = 0; i < this.cols.length; i++) {
          if (this.checkFeature(this.cols[i].feature, 'edit')) {
            this.$set(this.temp, this.cols[i].name, '')
          }
          if (this.checkFeature(this.cols[i].feature, 'search')) {
            if (this.listQuery[this.cols[i].name] === undefined) {
              this.$set(this.listQuery, this.cols[i].name, '')
            }
          }
          this.tHeader[i] = this.cols[i].name
          this.filterVal[i] = this.cols[i].name
        }

        this.listLoading = false
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.render()
    },
    checkFeature(feature, name) {
      return (feature & (1 << this.features[name])) > 0
    },
    sortChange(data) {
      const { prop, order } = data
      if (prop == null) {
        this.listQuery.sort = ''
      } else {
        if (order === 'descending') {
          this.listQuery.sort = '-' + prop
        } else {
          this.listQuery.sort = prop
        }
      }
      this.handleFilter()
    },
    resetTemp() {
      for (let i = 0; i < this.cols.length; i++) {
        this.$set(this.temp, this.cols[i].name, '')
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          createRow(this.listQuery.env_id, this.listQuery.collect, this.temp).then(() => {
            this.dialogFormVisible = false
            this.render()
            this.$notify({
              title: 'Success',
              message: 'Create Success',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row.values) // copy obj
      this.tempId = row.id
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          console.log(this.tempId)
          updateRow(this.listQuery.env_id, this.listQuery.collect, this.tempId, tempData).then(() => {
            this.dialogFormVisible = false
            this.render()
            this.$notify({
              title: 'Success',
              message: 'Update Success',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleConfirm() {
      if (this.dialogConfirmStatus === 'push') {
        this.pushData()
      }
      if (this.dialogConfirmStatus === 'delete') {
        this.deleteData()
      }
      if (this.dialogConfirmStatus === 'merge') {
        this.mergeData()
      }
    },
    handleDelete(row) {
      this.deleteRow = row
      this.dialogConfirmVisible = true
      this.dialogConfirmStatus = 'delete'
    },
    deleteData() {
      deleteRow(this.listQuery.env_id, this.listQuery.collect, this.deleteRow.id).then(() => {
        this.dialogConfirmVisible = false
        this.handleFilter()
        this.$notify({
          title: 'Success',
          message: 'Delete Success',
          type: 'success',
          duration: 2000
        })
      })
    },
    beforeUpload(file) {
      const isLt1M = file.size / 1024 / 1024 < 10

      if (isLt1M) {
        return true
      }

      this.$message({
        message: 'Please do not upload files larger than 10m in size.',
        type: 'warning'
      })
      return false
    },
    handleSuccess({ results, header }) {
      this.listLoading = true
      uploadData(this.listQuery.env_id, this.listQuery.collect, results).then(() => {
        this.listLoading = false
        this.handleFilter()
        this.$notify({
          title: 'Success',
          message: 'Upload Success',
          type: 'success',
          duration: 2000
        })
      })
    },
    handlePush() {
      this.dialogConfirmVisible = true
      this.dialogConfirmStatus = 'push'
    },
    pushData() {
      this.dialogConfirmVisible = false
      pushConfig(this.listQuery.env_id, this.pushServer, this.listQuery.collect).then(response => {
        this.dialogConfirmVisible = false
        const results = response.data.split(';')
        let timeout = 0
        const that = this
        for (let i = 0; i < results.length; i++) {
          const result = results[i]
          if (result === '') {
            continue
          }
          if (result.indexOf('SUCCESS') >= 0) {
            setTimeout(
              function() {
                that.$notify({
                  title: 'SUCCESS',
                  message: results[i],
                  type: 'success',
                  duration: 2000
                })
              }, timeout)
          } else {
            setTimeout(
              function() {
                that.$notify({
                  title: 'FAIL',
                  message: results[i],
                  type: 'warning',
                  duration: 2000
                })
              }, timeout)
          }

          timeout += 2000
        }
      })
    },
    handleMerge() {
      mergeConfigInfo(this.listQuery.env_id, this.listQuery.collect).then(response => {
        this.activities = response.data.activities
        this.merge_error = response.data.merge_error
        if (this.merge_error !== '') {
          this.dialogConfirmButton = false
        } else {
          this.dialogConfirmButton = true
        }
        this.dialogConfirmVisible = true
        this.dialogConfirmStatus = 'merge'
      })
    },
    mergeData() {
      this.dialogConfirmVisible = false
      mergeConfig(this.listQuery.env_id, this.listQuery.collect).then(() => {
        this.$notify({
          title: 'Success',
          message: 'Merge Success',
          type: 'success',
          duration: 2000
        })
      })
    }
  }
}
</script>
