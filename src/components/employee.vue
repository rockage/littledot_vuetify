<template>
  <div class="text-caption">
  Hello Products!
  </div>
</template>
<script>
export default {
  name: 'orders',
  data() {
    return {
    }
  },
  computed: {
    aMonthAgo: function () {
      let d = new Date()
      let str = ''
      str += d.getFullYear() + '-'
      str += d.getMonth() + '-'
      str += d.getDate()
      return str
    },
    Now: function () {
      let d = new Date()
      let str = ''
      str += d.getFullYear() + '-'
      str += d.getMonth() + 1 + '-'
      str += d.getDate()
      return str
    },
    viewModel: {
      // getter
      get: function () {
        return this.view_model
      },
      // setter
      set: function (newValue) {
        this.view_model = newValue
        this.initView()
      },
    },
  },
  methods: {
    TEST1: function () {},

    TEST2: function () {
      localStorage.setItem('state_list', '')
    },
    TEST3: function () {
    },
  },
  mounted: function () {
  },
}
</script>
<style>
.box {
  display: flex;
}

/* 表格专属CSS来自DivTable.com */
.divTable {
  display: table;
  width: 100%;
}

.divTableRow {
  display: table-row;
}

.divTableRow:hover {
  background-color: #eeeeee;
}

.divTableHeading {
  background-color: #eee;
  display: table-header-group;
}

.divTableCell,
.divTableHead {
  border: 1px solid #e0e0e0;
  display: table-cell;
  padding: 3px 10px;
  vertical-align: middle;
}

.center {
  text-align: center;
}

.nowrap {
  white-space: nowrap;
}

.divTableHeading {
  background-color: #eee;
  display: table-header-group;
  font-weight: bold;
}

.divTableFoot {
  background-color: #eee;
  display: table-footer-group;
  font-weight: bold;
}

.divTableBody {
  display: table-row-group;
}
</style>
