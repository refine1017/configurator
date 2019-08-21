<template>
  <div class="app-container">
    <div class="filter-container">

      <el-select v-model="listQuery.project_id" @change="handleChangeProject">
        <el-option v-for="item in projects" :key="item.id" :label="item.name" :value="item.id" />
      </el-select>

      <template v-for="(col ,index) in cols">
        <el-input v-if="checkFeature(col.feature, 'search')" v-model="listQuery[col.name]" :placeholder="col.name" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      </template>

      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        {{ $t('table.search') }}
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        {{ $t('table.add') }}
      </el-button>
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
        <el-table-column :prop="col.name" :sortable="checkFeature(col.feature, 'sort')" :label="col.name" align="center">
          <template slot-scope="scope">
            <span>{{ scope.row[col.name] }}</span>
          </template>
        </el-table-column>
      </template>

      <el-table-column :label="$t('table.actions')" align="center" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <router-link v-if="row.format === 'table'" :to="row.EditUrl">
            <el-button type="success" size="mini">
              Fields
            </el-button>
          </router-link>
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

    <el-dialog :title="textMap['confirm']" :visible.sync="dialogConfirmVisible">
      <el-alert
        :title="textMap['confirmMsg']"
        type="warning"
        effect="dark"
      />
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogConfirmVisible = false">
          {{ $t('table.cancel') }}
        </el-button>
        <el-button type="primary" @click="deleteData()">
          {{ $t('table.confirm') }}
        </el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { fetchList, createRow, updateRow, deleteRow } from '@/api/server'
import waves from '@/directive/waves' // waves directive
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import store from '@/store'
import TableForm from '@/components/TableForm'

export default {
  name: 'ComplexTable',
  components: { Pagination, TableForm },
  directives: { waves },
  filters: {},
  data() {
    return {
      projects: store.getters.userProjects,
      tableKey: 0,
      cols: null,
      list: null,
      total: 0,
      exportUrl: '',
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        project_id: store.getters.userSelectProject.id,
      },
      temp: {},
      tempId: -1,
      features: store.getters.features,
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: 'Edit',
        create: 'Create',
        confirm: 'Confirm',
        confirmMsg: 'Are you sure?'
      },
      dialogConfirmVisible: false,
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
    handleChangeProject() {
      for (let i = 0; i < this.projects.length; i++) {
        if (this.projects[i].id === this.listQuery.projectId) {
          this.environments = this.projects[i].envs
        }
      }
      if (this.environments.length > 0) {
        this.listQuery.env_id = this.environments[0].id
      }
      this.handleFilter()
    },
    handleFilter() {
      this.listQuery.page = 1
      this.render()
    },
    checkFeature(feature, name) {
      return (feature & (1 << this.features[name])) > 0
    },
    sortChange(data) {
    },
    resetTemp() {
      for (let i = 0; i < this.cols.length; i++) {
        if (this.checkFeature(this.cols[i].feature, 'edit')) {
          this.$set(this.temp, this.cols[i].name, '')
        }
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
          createRow(this.listQuery.project_id, this.temp).then(() => {
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
      const that = this
      Object.keys(row).forEach(function(key) {
        if (that.temp[key] !== undefined) {
          that.temp[key] = row[key]
        }
      })
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
          updateRow(this.tempId, tempData).then(() => {
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
    handleDelete(row) {
      this.deleteRow = row
      this.dialogConfirmVisible = true
    },
    deleteData() {
      deleteRow(this.deleteRow.id).then(() => {
        this.dialogConfirmVisible = false
        this.handleFilter()
        this.$notify({
          title: 'Success',
          message: 'Delete Success',
          type: 'success',
          duration: 2000
        })
      })
    }
  }
}
</script>
