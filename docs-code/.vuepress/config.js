module.exports = {
    title: 'mug Documentation',
    description: 'Documentation for mug which lets you create AWS Lambda for golang projects.',
    base: '/mug/',
    markdown: {
        toc: {
            includeLevel: [ 1, 2 ],
        },
        lineNumbers: true,
    },
    themeConfig: {
        nav: [
            { text: 'Guide', link: '/'},
            { text: 'Command Reference', link: '/commands/'},
        ],
        sidebar: {
            '/commands/': [
                {
                    title: 'Commands Reference',
                    collapsable: false,
                    children: [
                        'mug',
                        'mug_create',
                        'mug_add',
                        'mug_debug',
                        'mug_deploy',
                        'mug_remove',
                    ],
                }, 
            ],
            '/': [
                {
                    title: 'Guide',
                    collapsable: false,
                    children: [
                        '',
                        'getting-started',
                        'add',
                        'debug',
                        'deploy',
                    ],
                }, 
            ],
            
        }
    }
}