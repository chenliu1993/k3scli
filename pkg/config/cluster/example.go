package config

// A default cluster config

var (
	// DefaultClusterConfig :cluster config
	DefaultClusterConfig = Cluster{
		Name: DefaultClusterName,
		Nodes: []Node{
			{
				Name:  "worker0",
				Label: "worker",
				Ports: []Port{
					{
						Hostport: "8900",
						Port:     "8900",
					},
				},
			},
			{
				Name:  "server",
				Label: "server",
				Ports: []Port{
					{
						Hostport: "8901",
						Port:     "8901",
					},
					{Hostport: "6445",
						Port: "6443",
					},
				},
			},
			{
				Name:  "worker1",
				Label: "worker",
				Ports: []Port{
					{
						Hostport: "8902",
						Port:     "8902",
					},
				},
			},
		},
	}
)
