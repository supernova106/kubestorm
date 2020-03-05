<template>
  <div>
    <div>
      <b-container>
        <h2>Welcome to {{ msg }}</h2>
        <b-button variant="outline-primary" v-on:click="fetchData">List EKS Clusters</b-button>
        <div v-if="loading">
          <b-spinner small variant="primary" label="Spinning"></b-spinner>
        </div>
      </b-container>
    </div>
    <div>
      <b-tabs pills card align="center">
        <b-container v-for="result in results" :key="result.AccountID">
          <b-tab active>
            <template v-slot:title>{{ result.Metadata.Name }}</template>
            <p>EKS Count: {{ result.Clusters.length }}</p>
            <p>Environment: {{ result.Metadata.Environment }}</p>
            <p>AccountID: {{ result.AccountID }}</p>
            <b-table
              striped
              hover
              :items="result.Clusters"
              :fields="cluster_fields"
              :sort-by.sync="sortBy"
              :sort-desc.sync="sortDesc"
            ></b-table>
          </b-tab>
          <!-- <h3>{{ result.Metadata.Name }}</h3> -->
        </b-container>
      </b-tabs>
    </div>
  </div>
</template>

<script>
import axios from "axios";
export default {
  name: "Eks",
  props: {
    msg: String
  },
  data() {
    return {
      sortBy: "Version",
      sortDesc: true,
      cluster_fields: [
        {
          key: "Name",
          sortable: true
        },
        {
          key: "CreatedAt",
          sortable: true
        },
        {
          key: "Status",
          sortable: true
        },
        {
          key: "VpcId",
          sortable: true
        },
        {
          key: "Region",
          sortable: true
        },
        {
          key: "Version",
          sortable: true
        }
      ],
      results: [],
      loading: false
    };
  },
  methods: {
    fetchData: function() {
      {
        this.loading = true;
        axios.get(process.env.VUE_APP_KUBESTORM_URL).then(
          response => {
            this.loading = false;
            this.results = response.data;
          },
          error => {
            this.loading = false;
          }
        );
      }
    }
  }
};
</script>