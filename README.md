# mlu metrics
## example
```
# HELP HELP cnmon_mlu_temp GPU temperature (in C).
# TYPE cnmon_mlu_temp guage
cnmon_mlu_temp{mlu="0",uuid="011812120333"} 29
# HELP cnmon_mlu_board_utilization MLU utilization (in %).
# TYPE cnmon_mlu_board_utilization guage
cnmon_mlu_board_utilization{mlu="0",uuid="011812120333"} 0
# HELP cnmon_physical_memory_used/free MLU utilization (in MB).
# TYPE cnmon_physical_memory_used/free guage
cnmon_physical_memory_used{mlu="0",uuid="011812120333"} 4192
cnmon_physical_memory_free{mlu="0",uuid="011812120333"} 4000
# HELP cnmon_virtual_memory_used/free MLU utilization (in MB).
# TYPE cnmon_virtual_memory_used/free guage
cnmon_virtual_memory_used{mlu="0,",uuid="011812120333,"} 16643
cnmon_virtual_memory_free{mlu="0,",uuid="011812120333,"} 16125
```