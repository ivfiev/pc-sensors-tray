sleep 3
~/GolandProjects/pc-sensors-tray/pc-sensors-tray nvme-pci-0400 Composite.temp1_input "Sensor 1.temp2_input" "Sensor 2.temp3_input" &
sleep 0.5
~/GolandProjects/pc-sensors-tray/pc-sensors-tray amdgpu-pci-0300 edge.temp1_input junction.temp2_input mem.temp3_input fan1.fan1_input PPT.power1_average &
sleep 0.5
~/GolandProjects/pc-sensors-tray/pc-sensors-tray k10temp-pci-00c3 Tctl.temp1_input Tccd1.temp3_input Tccd2.temp4_input &
sleep 0.5
~/GolandProjects/pc-sensors-tray/pc-sensors-tray cpu-freq 3000 5600 &