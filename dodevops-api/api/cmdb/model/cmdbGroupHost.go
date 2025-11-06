package model

// 分组与主机关联DTO
type CmdbGroupHostDto struct {
	ID       uint               `json:"id"`       // 分组ID
	ParentID uint               `json:"parentId"` // 父分组ID
	Name     string             `json:"name"`     // 分组名称
	Hosts    []SimpleHostVo     `json:"hosts"`    // 主机列表(简化版)
	Children []CmdbGroupHostDto `json:"children"` // 子分组
}

// 简化版主机信息
type SimpleHostVo struct {
	ID       uint   `json:"id"`       // 主机ID
	HostName string `json:"hostName"` // 主机名称
	SSHIP    string `json:"sshIp"`    // SSH连接IP
	SSHName  string `json:"sshName"`  // SSH用户名
	SSHPort  int    `json:"sshPort"`  // SSH端口
}

// 构建分组树与主机关联结构
func BuildCmdbGroupHostTree(groups []CmdbGroup, hosts []CmdbHostVo) []CmdbGroupHostDto {
	groupMap := make(map[uint]CmdbGroupHostDto)
	hostsByGroup := make(map[uint][]SimpleHostVo)

	// 按分组ID组织主机并转换为简化版
	for _, host := range hosts {
		hostsByGroup[host.GroupID] = append(hostsByGroup[host.GroupID], SimpleHostVo{
			ID:       host.ID,
			HostName: host.HostName,
			SSHIP:    host.SSHIP,
			SSHName:  host.SSHName,
			SSHPort:  host.SSHPort,
		})
	}

	// 创建分组映射
	for _, group := range groups {
		groupMap[group.ID] = CmdbGroupHostDto{
			ID:       group.ID,
			ParentID: group.ParentID,
			Name:     group.Name,
			Hosts:    hostsByGroup[group.ID],
		}
	}

	// 构建树形结构
	var tree []CmdbGroupHostDto
	for _, group := range groups {
		if group.ParentID == 0 {
			tree = append(tree, buildGroupHostSubTree(groupMap[group.ID], groupMap))
		}
	}
	return tree
}

// 递归构建子树
func buildGroupHostSubTree(group CmdbGroupHostDto, groupMap map[uint]CmdbGroupHostDto) CmdbGroupHostDto {
	for _, child := range groupMap {
		if child.ParentID == group.ID {
			group.Children = append(group.Children, buildGroupHostSubTree(child, groupMap))
		}
	}
	return group
}
