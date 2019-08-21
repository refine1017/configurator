<template>
  <div class="app-container">
    <div class="filter-container">

      <template v-for="(col ,index) in cols">
        <el-input v-if="checkFeature(col.feature, 'search')" v-model="listQuery[col.name]" :placeholder="col.name" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      </template>

      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        {{ $t('table.search') }}
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
        <el-table-column v-if="checkFeature(col.feature, 'disableShow') === false" :prop="col.name" :sortable="checkFeature(col.feature, 'sort') ? 'custom' : false" :label="col.name" align="center">
          <template slot-scope="scope">
            <span>{{ scope.row[col.name] }}</span>
          </template>
        </el-table-column>
      </template>

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
import { fetchLogList } from '@/api/admin'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
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
      tableKey: 0,
      cols: null,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        sort: '-time'
      },
      temp: {},
      tempId: '',
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
      fetchLogList(this.listQuery).then(response => {
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
        if (this.checkFeature(this.cols[i].feature, 'edit')) {
          this.$set(this.temp, this.cols[i].name, '')
        }
      }
    },
    handleDownload() {
      this.downloadLoading = true
        import('@/vendor/Export2Excel').then(excel => {
          const newList = []
          for (let i = 0; i < this.list.length; i++) {
            newList.push(this.list[i])
          }

          const data = this.formatJson(this.filterVal, newList)
          excel.export_json_to_excel({
            header: this.tHeader,
            data,
            filename: this.listQuery.collect
          })
          this.downloadLoading = false
        })
    },
    formatJson(filterVal, jsonData) {
      return jsonData.map(v => filterVal.map(j => {
        if (j === 'timestamp') {
          return parseTime(v[j])
        } else {
          return v[j]
        }
      }))
    }
  }
}
</script>
